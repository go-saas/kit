package service

import (
	_ "embed"
	"github.com/flowchartsman/swaggerui"
	"github.com/go-chi/chi/v5"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authz"
	"github.com/goxiaoy/go-saas-kit/pkg/blob"
	kconf "github.com/goxiaoy/go-saas-kit/pkg/conf"
	"github.com/goxiaoy/go-saas-kit/pkg/job"
	"github.com/goxiaoy/go-saas-kit/pkg/server"
	"github.com/goxiaoy/go-saas-kit/sys/api"
	v1 "github.com/goxiaoy/go-saas-kit/sys/api/menu/v1"
	"github.com/hibiken/asynq"
	"net/http"
)

//go:embed openapi/api.swagger.json
var spec []byte

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewHttpServerRegister, NewGrpcServerRegister,
	NewMenuService, wire.Bind(new(v1.MenuServiceServer), new(*MenuService)))

type HttpServerRegister server.HttpServiceRegister
type GrpcServerRegister server.GrpcServiceRegister

func NewHttpServerRegister(
	menu *MenuService,
	authzSrv authz.Service,
	errEncoder khttp.EncodeErrorFunc,
	factory blob.Factory,
	dataCfg *kconf.Data,
	opt asynq.RedisConnOpt,
) HttpServerRegister {
	return server.HttpServiceRegisterFunc(func(srv *khttp.Server, middleware middleware.Middleware) {
		server.HandleBlobs("", dataCfg.Blobs, srv, factory)
		v1.RegisterMenuServiceHTTPServer(srv, menu)

		router := chi.NewRouter()
		router.Use(
			server.MiddlewareConvert(errEncoder, middleware))

		const apiPrefix = "/v1/sys/dev/swagger"

		router.Handle(apiPrefix+"*", http.StripPrefix(apiPrefix, server.AuthzGuardian(
			authzSrv, authz.RequirementList{
				authz.NewRequirement(authz.NewEntityResource("dev", "sys"), authz.AnyAction),
			}, errEncoder, swaggerui.Handler(spec),
		)))
		const asynqPrefix = "/v1/sys/asynqmon"
		router.Handle(asynqPrefix+"*", server.AuthzGuardian(
			authzSrv, authz.RequirementList{
				authz.NewRequirement(authz.NewEntityResource(api.ResourceDevJob, "*"), authz.AnyAction),
			}, errEncoder, job.NewUi(asynqPrefix, opt),
		))
		srv.HandlePrefix(apiPrefix, router)
		srv.HandlePrefix(asynqPrefix, router)

	})
}

func NewGrpcServerRegister(menu *MenuService) GrpcServerRegister {
	return server.GrpcServiceRegisterFunc(func(srv *grpc.Server, middleware middleware.Middleware) {
		v1.RegisterMenuServiceServer(srv, menu)
	})
}
