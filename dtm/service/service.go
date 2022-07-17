package service

import (
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	v1 "github.com/go-saas/kit/dtm/api/dtm/v1"
	"github.com/go-saas/kit/pkg/server"
)

func NewHttpServerRegister(msg *MsgService) server.HttpServiceRegister {
	return server.HttpServiceRegisterFunc(func(server *http.Server, middleware ...middleware.Middleware) {
		v1.RegisterMsgServiceHTTPServer(server, msg)
	})
}
func NewGrpcServerRegister(msg *MsgService) server.GrpcServiceRegister {
	return server.GrpcServiceRegisterFunc(func(server *grpc.Server, middleware ...middleware.Middleware) {
		v1.RegisterMsgServiceServer(server, msg)
	})
}
