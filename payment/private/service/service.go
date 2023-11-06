package service

import (
	_ "embed"
	"github.com/flowchartsman/swaggerui"
	"github.com/go-chi/chi/v5"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	v13 "github.com/go-saas/kit/payment/api/checkout/v1"
	v1 "github.com/go-saas/kit/payment/api/gateway/v1"
	v12 "github.com/go-saas/kit/payment/api/subscription/v1"
	kitdi "github.com/go-saas/kit/pkg/di"
	kitgrpc "github.com/go-saas/kit/pkg/server/grpc"
	kithttp "github.com/go-saas/kit/pkg/server/http"
	"github.com/go-saas/kit/pkg/stripe"
	"net/http"
)

//go:embed openapi/api.swagger.json
var spec []byte

// ProviderSet is service providers.
var ProviderSet = kitdi.NewSet(
	NewGrpcServerRegister,
	NewHttpServerRegister,
	NewPaymentService,
	NewSubscriptionService,
	NewCheckoutService,
	stripe.ProviderSet,
)

func NewHttpServerRegister(
	resEncoder khttp.EncodeResponseFunc,
	errEncoder khttp.EncodeErrorFunc,
	paymentSrv *PaymentService,
	subscription *SubscriptionService,
	checkout *CheckoutService) kithttp.ServiceRegister {
	return kithttp.ServiceRegisterFunc(func(srv *khttp.Server, middleware ...middleware.Middleware) {

		v1.RegisterPaymentGatewayServiceHTTPServer(srv, paymentSrv)
		v1.RegisterStripePaymentGatewayServiceHTTPServer(srv, paymentSrv)

		v12.RegisterSubscriptionServiceHTTPServer(srv, subscription)
		v13.RegisterCheckoutServiceHTTPServer(srv, checkout)

		swaggerRouter := chi.NewRouter()
		swaggerRouter.Use(
			kithttp.MiddlewareConvert(errEncoder, middleware...))
		const apiPrefix = "/v1/payment/dev/swagger"
		swaggerRouter.Handle(apiPrefix+"*", http.StripPrefix(apiPrefix, swaggerui.Handler(spec)))
		srv.HandlePrefix(apiPrefix, swaggerRouter)
	})
}

func NewGrpcServerRegister(
	paymentSrv *PaymentService, subscription *SubscriptionService, checkout *CheckoutService) kitgrpc.ServiceRegister {
	return kitgrpc.ServiceRegisterFunc(func(srv *grpc.Server, middleware ...middleware.Middleware) {
		v1.RegisterPaymentGatewayServiceServer(srv, paymentSrv)
		v1.RegisterStripePaymentGatewayServiceServer(srv, paymentSrv)
		v12.RegisterSubscriptionServiceServer(srv, subscription)
		v12.RegisterSubscriptionInternalServiceServer(srv, subscription)
		v13.RegisterCheckoutServiceServer(srv, checkout)
	})
}
