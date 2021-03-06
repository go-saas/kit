package service

import (
	_ "embed"
	"github.com/flowchartsman/swaggerui"
	"github.com/go-chi/chi/v5"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-saas/kit/pkg/blob"
	kconf "github.com/go-saas/kit/pkg/conf"
	kitdi "github.com/go-saas/kit/pkg/di"
	"github.com/go-saas/kit/pkg/server"
	v1 "github.com/go-saas/kit/saas/api/tenant/v1"
	"github.com/goava/di"
	"net/http"
)

//go:embed openapi/api.swagger.json
var spec []byte

// ProviderSet is service providers.
var ProviderSet = kitdi.NewSet(NewHttpServerRegister, NewGrpcServerRegister,
	kitdi.NewProvider(NewTenantService, di.As(new(v1.TenantServiceServer))),
	kitdi.NewProvider(NewTenantInternalService, di.As(new(v1.TenantInternalServiceServer))),
)

func NewHttpServerRegister(
	tenant *TenantService,
	factory blob.Factory,
	errEncoder khttp.EncodeErrorFunc,
	tenantInternal *TenantInternalService,
	dataCfg *kconf.Data,
) server.HttpServiceRegister {
	return server.HttpServiceRegisterFunc(func(srv *khttp.Server, middleware ...middleware.Middleware) {
		route := srv.Route("/")

		route.POST("/v1/saas/tenant/logo", tenant.UpdateLogo)
		server.HandleBlobs("", dataCfg.Blobs, srv, factory)

		v1.RegisterTenantServiceHTTPServer(srv, tenant)

		router := chi.NewRouter()
		//global filter
		router.Use(
			server.MiddlewareConvert(errEncoder, middleware...))
		const apiPrefix = "/v1/saas/dev/swagger"
		router.Handle(apiPrefix+"*", http.StripPrefix(apiPrefix, swaggerui.Handler(spec)))
		srv.HandlePrefix(apiPrefix, router)
	})
}

func NewGrpcServerRegister(tenant *TenantService, tenantInternal *TenantInternalService) server.GrpcServiceRegister {
	return server.GrpcServiceRegisterFunc(func(srv *grpc.Server, middleware ...middleware.Middleware) {
		v1.RegisterTenantInternalServiceServer(srv, tenantInternal)
		v1.RegisterTenantServiceServer(srv, tenant)
	})
}
