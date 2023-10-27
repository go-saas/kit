package service

import (
	_ "embed"
	"github.com/flowchartsman/swaggerui"
	"github.com/go-chi/chi/v5"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	kitdi "github.com/go-saas/kit/pkg/di"
	kitgrpc "github.com/go-saas/kit/pkg/server/grpc"
	kithttp "github.com/go-saas/kit/pkg/server/http"
	"github.com/go-saas/kit/pkg/stripe"
	v12 "github.com/go-saas/kit/product/api/product/v1"
	"net/http"
)

//go:embed openapi/api.swagger.json
var spec []byte

// ProviderSet is service providers.
var ProviderSet = kitdi.NewSet(
	stripe.ProviderSet,

	NewGrpcServerRegister,
	NewHttpServerRegister,
	NewUploadService,
	NewProductService,

	NewProductUpdatedTaskHandler,
)

func NewHttpServerRegister(
	resEncoder khttp.EncodeResponseFunc,
	errEncoder khttp.EncodeErrorFunc,
	product *ProductService) kithttp.ServiceRegister {
	return kithttp.ServiceRegisterFunc(func(srv *khttp.Server, middleware ...middleware.Middleware) {
		v12.RegisterProductServiceHTTPServer(srv, product)

		route := srv.Route("/")
		route.POST("/v1/product/media", product.UploadMedias)

		swaggerRouter := chi.NewRouter()
		swaggerRouter.Use(
			kithttp.MiddlewareConvert(errEncoder, middleware...))
		const apiPrefix = "/v1/product/dev/swagger"
		swaggerRouter.Handle(apiPrefix+"*", http.StripPrefix(apiPrefix, swaggerui.Handler(spec)))
		srv.HandlePrefix(apiPrefix, swaggerRouter)
	})
}

func NewGrpcServerRegister(product *ProductService) kitgrpc.ServiceRegister {
	return kitgrpc.ServiceRegisterFunc(func(srv *grpc.Server, middleware ...middleware.Middleware) {
		v12.RegisterProductServiceServer(srv, product)

		v12.RegisterProductInternalServiceServer(srv, product)
	})
}
