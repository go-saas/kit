package main

import (
	"errors"
	"fmt"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	klog "github.com/go-kratos/kratos/v2/log"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/gateway/apisix/cmd/go-runner/plugins"
	"github.com/goxiaoy/go-saas-kit/pkg/api"
	"github.com/goxiaoy/go-saas-kit/pkg/authn/jwt"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	v1 "github.com/goxiaoy/go-saas-kit/saas/api/tenant/v1"
	uremote "github.com/goxiaoy/go-saas-kit/user/remote"
	"github.com/goxiaoy/go-saas/common"
	shttp "github.com/goxiaoy/go-saas/common/http"
	sapisix "github.com/goxiaoy/go-saas/gateway/apisix"
	"net/http"
)

type App struct {
	tenantStore  common.TenantStore
	tokenizer    jwt.Tokenizer
	tokenManager api.TokenManager
	services     *conf.Services
	clientName   api.ClientName
	logger       klog.Logger
	tenantCfg    *shttp.WebMultiTenancyOption
	security     *conf.Security
	userTenant   *uremote.UserTenantContributor
}

func newApp(tenantStore common.TenantStore,
	userTenant *uremote.UserTenantContributor,
	t jwt.Tokenizer,
	tmr api.TokenManager,
	services *conf.Services,
	clientName api.ClientName,
	tenantCfg *shttp.WebMultiTenancyOption,
	security *conf.Security,
	logger klog.Logger) (*App, error) {
	ret := &App{tenantStore: tenantStore,
		userTenant:   userTenant,
		tokenizer:    t,
		tokenManager: tmr,
		services:     services,
		clientName:   clientName,
		tenantCfg:    tenantCfg,
		security:     security,
		logger:       logger}
	return ret, ret.load()
}

func (a *App) load() error {
	sapisix.Init(a.tenantStore, fmt.Sprintf("%s%s", api.PrefixOrDefault(""), a.tenantCfg.TenantKey), func(err error, w http.ResponseWriter) {
		if errors.Is(err, common.ErrTenantNotFound) {
			err = v1.ErrorTenantNotFound("")
		}
		//use error codec
		fr := kerrors.FromError(err)
		w.WriteHeader(int(fr.Code))
		khttp.DefaultErrorEncoder(w, &http.Request{}, err)
	})
	if err := plugins.Init(a.tokenizer, a.tokenManager, a.clientName, a.services, a.security, a.userTenant, a.logger); err != nil {
		return err
	}
	return nil
}

func NewSelfClientOption(logger log.Logger) *api.Option {
	return api.NewOption("", false, api.NewUserContributor(logger), api.NewClientContributor(true, logger))
}

var ProviderSet = wire.NewSet(api.NewInMemoryTokenManager, NewSelfClientOption,
	wire.Bind(new(api.TokenManager), new(*api.InMemoryTokenManager)), shttp.NewDefaultWebMultiTenancyOption)
