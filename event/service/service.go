package service

import (
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	v12 "github.com/go-saas/kit/event/api/v1"
	kitgrpc "github.com/go-saas/kit/pkg/server/grpc"
	kithttp "github.com/go-saas/kit/pkg/server/http"
)

func NewHttpServerRegister(event *EventService) kithttp.ServiceRegister {
	return kithttp.ServiceRegisterFunc(func(server *http.Server, middleware ...middleware.Middleware) {

	})
}
func NewGrpcServerRegister(event *EventService) kitgrpc.ServiceRegister {
	return kitgrpc.ServiceRegisterFunc(func(server *grpc.Server, middleware ...middleware.Middleware) {
		v12.RegisterEventServiceServer(server, event)
	})
}
