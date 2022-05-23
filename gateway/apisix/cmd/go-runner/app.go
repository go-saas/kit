package main

import (
	"github.com/go-kratos/kratos/v2/log"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/gateway/apisix/cmd/go-runner/plugins"
	"github.com/goxiaoy/go-saas-kit/pkg/api"
	"github.com/goxiaoy/go-saas-kit/pkg/authn/jwt"
	"github.com/goxiaoy/go-saas-kit/pkg/authn/session"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authz"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	uremote "github.com/goxiaoy/go-saas-kit/user/remote"
	"github.com/goxiaoy/go-saas/common"
	shttp "github.com/goxiaoy/go-saas/common/http"
)

type App struct {
	tenantStore     common.TenantStore
	tokenizer       jwt.Tokenizer
	tokenManager    api.TokenManager
	services        *conf.Services
	clientName      api.ClientName
	logger          klog.Logger
	tenantCfg       *shttp.WebMultiTenancyOption
	security        *conf.Security
	userTenant      *uremote.UserTenantContributor
	refreshProvider session.RefreshTokenProvider
	authService     authz.Service
	subjectResolver authz.SubjectResolver
}

func newApp(tenantStore common.TenantStore,
	userTenant *uremote.UserTenantContributor,
	t jwt.Tokenizer,
	tmr api.TokenManager,
	services *conf.Services,
	clientName api.ClientName,
	tenantCfg *shttp.WebMultiTenancyOption,
	security *conf.Security,
	refreshProvider session.RefreshTokenProvider,
	authService authz.Service,
	subjectResolver authz.SubjectResolver,

	logger klog.Logger) (*App, error) {
	ret := &App{tenantStore: tenantStore,
		userTenant:      userTenant,
		tokenizer:       t,
		tokenManager:    tmr,
		services:        services,
		clientName:      clientName,
		tenantCfg:       tenantCfg,
		security:        security,
		refreshProvider: refreshProvider,
		authService:     authService,
		subjectResolver: subjectResolver,
		logger:          logger}
	return ret, ret.load()
}

func (a *App) load() error {
	if err := plugins.Init(
		a.tokenizer,
		a.tokenManager,
		a.tenantCfg,
		a.clientName,
		a.services,
		a.security,
		a.userTenant,
		a.tenantStore,
		a.refreshProvider,
		a.authService,
		a.subjectResolver,
		a.logger,
	); err != nil {
		return err
	}
	return nil
}

func NewSelfClientOption(logger log.Logger) *api.Option {
	return api.NewOption(
		false,
		api.NewSaasPropagator(logger),
		api.NewUserPropagator(logger),
		//do not propagate client
		api.NewClientPropagator(true, logger),
	).WithInsecure()
}

func NewAuthorizationOption() *authz.Option {
	return authz.NewAuthorizationOption()
}

var ProviderSet = wire.NewSet(api.NewInMemoryTokenManager, NewSelfClientOption,
	wire.Bind(new(api.TokenManager), new(*api.InMemoryTokenManager)), NewAuthorizationOption)
