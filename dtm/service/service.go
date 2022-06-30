package service

import (
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	v1 "github.com/go-saas/kit/dtm/api/dtm/v1"
	"github.com/go-saas/kit/pkg/server"
)

type HttpServerRegister server.HttpServiceRegister
type GrpcServerRegister server.GrpcServiceRegister

func NewHttpServerRegister(msg *MsgServiceService) HttpServerRegister {
	return server.HttpServiceRegisterFunc(func(server *http.Server, middleware ...middleware.Middleware) {
		v1.RegisterMsgServiceHTTPServer(server, msg)
	})
}
func NewGrpcServerRegister(msg *MsgServiceService) GrpcServerRegister {
	return server.GrpcServiceRegisterFunc(func(server *grpc.Server, middleware ...middleware.Middleware) {
		v1.RegisterMsgServiceServer(server, msg)
	})
}
