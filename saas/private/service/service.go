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
	v12 "github.com/go-saas/kit/saas/api/plan/v1"
	v1 "github.com/go-saas/kit/saas/api/tenant/v1"
	"github.com/go-saas/kit/saas/private/biz"
	"github.com/goava/di"
	"github.com/goxiaoy/vfs"
	"net/http"
)

//go:embed openapi/api.swagger.json
var spec []byte

// ProviderSet is service providers.
var ProviderSet = kitdi.NewSet(NewHttpServerRegister, NewGrpcServerRegister,
	kitdi.NewProvider(NewTenantService, di.As(new(v1.TenantServiceServer))),
	kitdi.NewProvider(NewTenantInternalService, di.As(new(v1.TenantInternalServiceServer))),
	kitdi.NewProvider(NewPlanService, di.As(new(v12.PlanServiceServer))),
)

func NewHttpServerRegister(
	tenant *TenantService,
	errEncoder khttp.EncodeErrorFunc,
	tenantInternal *TenantInternalService,
	plan *PlanService,
	blob vfs.Blob,
) kithttp.ServiceRegister {
	return kithttp.ServiceRegisterFunc(func(srv *khttp.Server, middleware ...middleware.Middleware) {
		route := srv.Route("/")

		route.POST("/v1/saas/tenant/logo", tenant.UpdateLogo)
		kithttp.MountBlob(srv, "", biz.TenantLogoPath, blob)

		v1.RegisterTenantServiceHTTPServer(srv, tenant)
		v12.RegisterPlanServiceHTTPServer(srv, plan)

		router := chi.NewRouter()
		//global filter
		router.Use(
			kithttp.MiddlewareConvert(errEncoder, middleware...))
		const apiPrefix = "/v1/saas/dev/swagger"
		router.Handle(apiPrefix+"*", http.StripPrefix(apiPrefix, swaggerui.Handler(spec)))
		srv.HandlePrefix(apiPrefix, router)
	})
}

func NewGrpcServerRegister(tenant *TenantService, tenantInternal *TenantInternalService, plan *PlanService) kitgrpc.ServiceRegister {
	return kitgrpc.ServiceRegisterFunc(func(srv *grpc.Server, middleware ...middleware.Middleware) {
		v1.RegisterTenantInternalServiceServer(srv, tenantInternal)
		v1.RegisterTenantServiceServer(srv, tenant)

		v12.RegisterPlanServiceServer(srv, plan)
	})
}
