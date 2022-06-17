package service

import (
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authz"
	"github.com/goxiaoy/go-saas-kit/pkg/blob"
	"github.com/goxiaoy/go-saas-kit/pkg/server"
	v1 "github.com/goxiaoy/go-saas-kit/saas/api/tenant/v1"
	conf2 "github.com/goxiaoy/go-saas-kit/saas/private/conf"
)

func NewAuthorizationOption() *authz.Option {
	return authz.NewAuthorizationOption()
}

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewHttpServerRegister, NewGrpcServerRegister, NewTenantService, NewAuthorizationOption)

type HttpServerRegister server.HttpServiceRegister
type GrpcServerRegister server.GrpcServiceRegister

func NewHttpServerRegister(tenant *TenantService, factory blob.Factory,
	dataCfg *conf2.Data) HttpServerRegister {
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
