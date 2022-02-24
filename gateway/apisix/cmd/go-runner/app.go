package main

import (
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/gateway/apisix/cmd/go-runner/plugins"
	"github.com/goxiaoy/go-saas-kit/pkg/api"
	"github.com/goxiaoy/go-saas-kit/pkg/authn/jwt"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	"github.com/goxiaoy/go-saas/common"
	sapisix "github.com/goxiaoy/go-saas/gateway/apisix"
)

type App struct {
	tenantStore  common.TenantStore
	tokenizer    jwt.Tokenizer
	tokenManager api.TokenManager
	services     *conf.Services
	clientName   api.ClientName
	ao           *api.Option
}

func newApp(tenantStore common.TenantStore, t jwt.Tokenizer, tmr api.TokenManager, services *conf.Services, clientName api.ClientName, ao *api.Option) (*App, error) {
	ret := &App{tenantStore: tenantStore, tokenizer: t, tokenManager: tmr, services: services, clientName: clientName, ao: ao}
	return ret, ret.load()
}

func (a *App) load() error {
	sapisix.InitTenantStore(a.tenantStore)
	if err := plugins.Init(a.tokenizer, a.tokenManager, a.clientName, a.services, a.ao); err != nil {
		return err
	}
	return nil
}

func NewOption() *api.Option {
	return api.NewOption("", true, api.NewUserContributor(), api.NewClientContributor())
}

var ProviderSet = wire.NewSet(api.NewInMemoryTokenManager, NewOption,
	wire.Bind(new(api.TokenManager), new(*api.InMemoryTokenManager)))
