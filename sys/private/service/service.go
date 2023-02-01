package service

import (
	_ "embed"
	"github.com/flowchartsman/swaggerui"
	"github.com/go-chi/chi/v5"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	oidcservice "github.com/go-saas/kit/oidc/service"
	"github.com/go-saas/kit/pkg/apisix"
	"github.com/go-saas/kit/pkg/authz/authz"
	kconf "github.com/go-saas/kit/pkg/conf"
	kitdi "github.com/go-saas/kit/pkg/di"
	"github.com/go-saas/kit/pkg/job"
	kitgrpc "github.com/go-saas/kit/pkg/server/grpc"
	kithttp "github.com/go-saas/kit/pkg/server/http"
	"github.com/go-saas/kit/sys/api"
	v12 "github.com/go-saas/kit/sys/api/locale/v1"
	v1 "github.com/go-saas/kit/sys/api/menu/v1"
	"github.com/go-saas/kit/sys/private/conf"
	"github.com/goava/di"
	"github.com/hibiken/asynq"
	"net/http"
)

//go:embed openapi/api.swagger.json
var spec []byte

// ProviderSet is service providers.
var ProviderSet = kitdi.NewSet(NewApisixOption, NewApisixAdminClient, apisix.NewWatchSyncAdmin, oidcservice.ProviderSet,
	NewHttpServerRegister, NewGrpcServerRegister,
	kitdi.NewProvider(NewMenuService, di.As(new(v1.MenuServiceServer))),
	kitdi.NewProvider(NewLocaleService, di.As(new(v12.LocaleServiceServer))),
)

func NewApisixAdminClient(cfg *conf.SysConf) (*apisix.AdminClient, error) {
	var endpoint, apikey string
	if cfg != nil {
		if cfg.Apisix != nil {
			endpoint = cfg.Apisix.Endpoint
			apikey = cfg.Apisix.ApiKey
		}
	}
	return apisix.NewAdminClient(endpoint, apikey)
}

func NewApisixOption(srvs *kconf.Services) *apisix.Option {
	ret := &apisix.Option{
		Services: nil,
		Timeout:  0,
	}
	if srvs != nil {
		if srvs.Services != nil {
			for k, _ := range srvs.Services {
				if k != "default" {
					ret.Services = append(ret.Services, k)
				}
			}
		}
	}
	return ret
}

func NewHttpServerRegister(
	menu *MenuService,
	locSrv *LocaleService,
	authzSrv authz.Service,
	errEncoder khttp.EncodeErrorFunc,
	opt asynq.RedisConnOpt,
) kithttp.ServiceRegister {
	return kithttp.ServiceRegisterFunc(func(srv *khttp.Server, middleware ...middleware.Middleware) {

		v1.RegisterMenuServiceHTTPServer(srv, menu)
		v12.RegisterLocaleServiceHTTPServer(srv, locSrv)

		router := chi.NewRouter()
		router.Use(
			kithttp.MiddlewareConvert(errEncoder, middleware...))

		const apiPrefix = "/v1/sys/dev/swagger"

		router.Handle(apiPrefix+"*", http.StripPrefix(apiPrefix, swaggerui.Handler(spec)))
		const asynqPrefix = "/v1/sys/asynqmon"
		router.Handle(asynqPrefix+"*", kithttp.AuthzGuardian(
			authzSrv, authz.RequirementList{
				authz.NewRequirement(authz.NewEntityResource(api.ResourceDevJob, "*"), authz.AnyAction),
			}, errEncoder, job.NewUi(asynqPrefix, opt),
		))
		srv.HandlePrefix(apiPrefix, router)
		srv.HandlePrefix(asynqPrefix, router)

	})
}

func NewGrpcServerRegister(
	menu *MenuService,
	locSrv *LocaleService) kitgrpc.ServiceRegister {
	return kitgrpc.ServiceRegisterFunc(func(srv *grpc.Server, middleware ...middleware.Middleware) {
		v1.RegisterMenuServiceServer(srv, menu)
		v12.RegisterLocaleServiceServer(srv, locSrv)
	})
}
