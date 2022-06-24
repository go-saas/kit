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
	"github.com/goxiaoy/go-saas-kit/pkg/server"
	v1 "github.com/goxiaoy/go-saas-kit/saas/api/tenant/v1"
	"net/http"
)

//go:embed openapi/api.swagger.json
var spec []byte

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewHttpServerRegister, NewGrpcServerRegister,
	NewTenantService, wire.Bind(new(v1.TenantServiceServer), new(*TenantService)),
	wire.Struct(new(TenantInternalService), "*"), wire.Bind(new(v1.TenantInternalServiceServer), new(*TenantInternalService)))

type HttpServerRegister server.HttpServiceRegister
type GrpcServerRegister server.GrpcServiceRegister

func NewHttpServerRegister(
	tenant *TenantService,
	factory blob.Factory,
	authzSrv authz.Service,
	errEncoder khttp.EncodeErrorFunc,
	tenantInternal *TenantInternalService,
	dataCfg *kconf.Data,
) HttpServerRegister {
	return server.HttpServiceRegisterFunc(func(srv *khttp.Server, middleware middleware.Middleware) {
		route := srv.Route("/")

		route.POST("/v1/saas/tenant/logo", tenant.UpdateLogo)
		server.HandleBlobs("", dataCfg.Blobs, srv, factory)

		v1.RegisterTenantServiceHTTPServer(srv, tenant)

		router := chi.NewRouter()
		//global filter
		router.Use(
			server.MiddlewareConvert(errEncoder, middleware))
		const apiPrefix = "/v1/saas/dev/swagger"
		swaggerRouter := router.Group(func(router chi.Router) {
			router.Handle(apiPrefix+"*", http.StripPrefix(apiPrefix, server.AuthzGuardian(
				authzSrv, authz.RequirementList{
					authz.NewRequirement(authz.NewEntityResource("dev", "saas"), authz.AnyAction),
				}, errEncoder, swaggerui.Handler(spec),
			)))
		})
		srv.HandlePrefix(apiPrefix, swaggerRouter)
	})
}

func NewGrpcServerRegister(tenant *TenantService, tenantInternal *TenantInternalService) GrpcServerRegister {
	return server.GrpcServiceRegisterFunc(func(srv *grpc.Server, middleware middleware.Middleware) {
		v1.RegisterTenantInternalServiceServer(srv, tenantInternal)
		v1.RegisterTenantServiceServer(srv, tenant)
	})
}
