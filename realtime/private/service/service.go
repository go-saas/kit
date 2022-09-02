package service

import (
	_ "embed"
	"github.com/centrifugal/centrifuge"
	"github.com/flowchartsman/swaggerui"
	"github.com/go-chi/chi/v5"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	kitdi "github.com/go-saas/kit/pkg/di"
	"github.com/go-saas/kit/pkg/server"
	v1 "github.com/go-saas/kit/realtime/api/notification/v1"
	"net/http"
)

//go:embed openapi/api.swagger.json
var spec []byte

// ProviderSet is service providers.
var ProviderSet = kitdi.NewSet(
	NewGrpcServerRegister,
	NewHttpServerRegister,
	NewCentrifugeRegister,
	NewCentrifugeNode,
	NewNotificationService,
	NewNotificationEventHandler,
)

func NewHttpServerRegister(
	errEncoder khttp.EncodeErrorFunc,
	notification *NotificationService) server.HttpServiceRegister {
	return server.HttpServiceRegisterFunc(func(srv *khttp.Server, middleware ...middleware.Middleware) {
		v1.RegisterNotificationServiceHTTPServer(srv, notification)

		swaggerRouter := chi.NewRouter()
		swaggerRouter.Use(
			server.MiddlewareConvert(errEncoder, middleware...))
		const apiPrefix = "/v1/server/dev/swagger"
		swaggerRouter.Handle(apiPrefix+"*", http.StripPrefix(apiPrefix, swaggerui.Handler(spec)))
	})
}

func NewGrpcServerRegister(notification *NotificationService) server.GrpcServiceRegister {
	return server.GrpcServiceRegisterFunc(func(srv *grpc.Server, middleware ...middleware.Middleware) {
		v1.RegisterNotificationServiceServer(srv, notification)
	})
}

func NewCentrifugeRegister(node *centrifuge.Node, errEncoder khttp.EncodeErrorFunc) server.HttpServiceRegister {
	return server.HttpServiceRegisterFunc(func(srv *khttp.Server, middleware ...middleware.Middleware) {
		websocketHandler := centrifuge.NewWebsocketHandler(node, centrifuge.WebsocketConfig{
			ReadBufferSize:     1024,
			UseWriteBufferPool: true,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		})
		r := chi.NewRouter()
		r.Use(
			server.MiddlewareConvert(errEncoder, middleware...))
		const apiPrefix = "/v1/realtime/connect/ws"
		r.Handle(apiPrefix+"*", http.StripPrefix(apiPrefix, auth(websocketHandler)))
		srv.HandlePrefix(apiPrefix, r)
	})
}
