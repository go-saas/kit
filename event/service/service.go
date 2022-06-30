package service

import (
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	v12 "github.com/go-saas/kit/event/api/v1"
	"github.com/go-saas/kit/pkg/server"
)

type HttpServerRegister server.HttpServiceRegister
type GrpcServerRegister server.GrpcServiceRegister

func NewHttpServerRegister(event *EventService) HttpServerRegister {
	return server.HttpServiceRegisterFunc(func(server *http.Server, middleware ...middleware.Middleware) {
		v12.RegisterEventServiceHTTPServer(server, event)
	})
}
func NewGrpcServerRegister(event *EventService) GrpcServerRegister {
	return server.GrpcServiceRegisterFunc(func(server *grpc.Server, middleware ...middleware.Middleware) {
		v12.RegisterEventServiceServer(server, event)
	})
}
