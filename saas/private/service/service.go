package service

import (
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/pkg/blob"
	kconf "github.com/goxiaoy/go-saas-kit/pkg/conf"
	"github.com/goxiaoy/go-saas-kit/pkg/server"
	v1 "github.com/goxiaoy/go-saas-kit/saas/api/tenant/v1"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewHttpServerRegister, NewGrpcServerRegister,
	NewTenantService, wire.Bind(new(v1.TenantServiceServer), new(*TenantService)),
	wire.Struct(new(TenantInternalService), "*"), wire.Bind(new(v1.TenantInternalServiceServer), new(*TenantInternalService)))

type HttpServerRegister server.HttpServiceRegister
type GrpcServerRegister server.GrpcServiceRegister

func NewHttpServerRegister(tenant *TenantService, factory blob.Factory, tenantInternal *TenantInternalService,
	dataCfg *kconf.Data) HttpServerRegister {
	return server.HttpServiceRegisterFunc(func(srv *http.Server, middleware middleware.Middleware) {
		route := srv.Route("/")

		route.POST("/v1/saas/tenant/logo", tenant.UpdateLogo)
		server.HandleBlobs("", dataCfg.Blobs, srv, factory)

		v1.RegisterTenantServiceHTTPServer(srv, tenant)
	})
}

func NewGrpcServerRegister(tenant *TenantService, tenantInternal *TenantInternalService) GrpcServerRegister {
	return server.GrpcServiceRegisterFunc(func(srv *grpc.Server, middleware middleware.Middleware) {
		v1.RegisterTenantInternalServiceServer(srv, tenantInternal)
		v1.RegisterTenantServiceServer(srv, tenant)
	})
}
