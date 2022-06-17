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
var ProviderSet = wire.NewSet(NewHttpServerRegister, NewGrpcServerRegister, NewTenantService, NewTenantServiceServer)

type HttpServerRegister server.HttpServiceRegister
type GrpcServerRegister server.GrpcServiceRegister

func NewHttpServerRegister(tenant *TenantService, factory blob.Factory,
	dataCfg *kconf.Data) HttpServerRegister {
	return server.HttpServiceRegisterFunc(func(srv *http.Server, middleware middleware.Middleware) {
		route := srv.Route("/")

		route.POST("/v1/saas/tenant/logo", tenant.UpdateLogo)
		server.HandleBlobs("", dataCfg.Blobs, srv, factory)
		v1.RegisterTenantServiceHTTPServer(srv, tenant)
	})
}

func NewGrpcServerRegister(tenant *TenantService) GrpcServerRegister {
	return server.GrpcServiceRegisterFunc(func(srv *grpc.Server, middleware middleware.Middleware) {
		v1.RegisterTenantServiceServer(srv, tenant)
	})
}
