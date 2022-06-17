package service

import (
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authz"
	"github.com/goxiaoy/go-saas-kit/pkg/blob"
	"github.com/goxiaoy/go-saas-kit/pkg/server"
	v1 "github.com/goxiaoy/go-saas-kit/sys/api/menu/v1"
	"github.com/goxiaoy/go-saas-kit/sys/private/conf"
)

func NewAuthorizationOption() *authz.Option {
	return authz.NewAuthorizationOption()
}

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewHttpServerRegister, NewGrpcServerRegister, NewAuthorizationOption, NewMenuService)

type HttpServerRegister server.HttpServiceRegister
type GrpcServerRegister server.GrpcServiceRegister

func NewHttpServerRegister(menu *MenuService, factory blob.Factory, dataCfg *conf.Data) HttpServerRegister {
	return server.HttpServiceRegisterFunc(func(srv *http.Server, middleware middleware.Middleware) {
		server.HandleBlobs("", dataCfg.Blobs, srv, factory)
		v1.RegisterMenuServiceHTTPServer(srv, menu)
	})
}

func NewGrpcServerRegister(menu *MenuService) GrpcServerRegister {
	return server.GrpcServiceRegisterFunc(func(srv *grpc.Server, middleware middleware.Middleware) {
		v1.RegisterMenuServiceServer(srv, menu)
	})
}
