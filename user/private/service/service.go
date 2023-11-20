package service

import (
	_ "embed"
	"github.com/flowchartsman/swaggerui"
	"github.com/go-chi/chi/v5"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	kitdi "github.com/go-saas/kit/pkg/di"
	"github.com/go-saas/kit/pkg/idp"
	kitgrpc "github.com/go-saas/kit/pkg/server/grpc"
	kithttp "github.com/go-saas/kit/pkg/server/http"
	"github.com/go-saas/kit/pkg/stripe"
	"github.com/go-saas/kit/user/api"
	v13 "github.com/go-saas/kit/user/api/account/v1"
	v14 "github.com/go-saas/kit/user/api/auth/v1"
	v15 "github.com/go-saas/kit/user/api/permission/v1"
	v1 "github.com/go-saas/kit/user/api/role/v1"
	v12 "github.com/go-saas/kit/user/api/user/v1"
	"github.com/go-saas/kit/user/private/biz"
	uhttp "github.com/go-saas/kit/user/private/service/http"
	"github.com/goava/di"
	"github.com/goxiaoy/vfs"
	"net/http"
)

//go:embed openapi/api.swagger.json
var spec []byte

// ProviderSet is service providers.
var ProviderSet = kitdi.NewSet(
	NewGrpcServerRegister,
	NewHttpServerRegister,
	NewUserRoleContrib,
	//idp
	idp.NewWeChat,
	stripe.ProviderSet,

	kitdi.NewProvider(NewUserService, di.As(new(v12.UserServiceServer), new(v12.UserAdminServiceServer))),
	kitdi.NewProvider(NewUserInternalService, di.As(new(v12.UserInternalServiceServer))),

	kitdi.NewProvider(NewAccountService, di.As(new(v13.AccountServer))),

	kitdi.NewProvider(NewAuthService, di.As(new(v14.AuthServer))),

	kitdi.NewProvider(NewRoleServiceService, di.As(new(v1.RoleServiceServer))),
	kitdi.NewProvider(NewPermissionService, di.As(new(v15.PermissionServiceServer))),

	kitdi.NewProvider(NewWeChatAuthService, di.As(new(v14.WeChatAuthServiceServer))),
	api.NewUserTenantContrib,
	api.NewRefreshProvider,
	uhttp.NewAuth)

func NewHttpServerRegister(
	user *UserService,
	userInternal *UserInternalService,
	resEncoder khttp.EncodeResponseFunc,
	errEncoder khttp.EncodeErrorFunc,
	account *AccountService,
	auth *AuthService,
	role *RoleService,
	permission *PermissionService,
	authHttp *uhttp.Auth,
	weChatAuth *WeChatAuthService,
	vfs vfs.Blob) kithttp.ServiceRegister {
	return kithttp.ServiceRegisterFunc(func(srv *khttp.Server, middleware ...middleware.Middleware) {

		router := chi.NewRouter()

		//global filter
		router.Use(
			kithttp.MiddlewareConvert(errEncoder, middleware...))

		router.Get("/login", kithttp.HandlerWrap(resEncoder, authHttp.LoginGet))
		router.Post("/login", kithttp.HandlerWrap(resEncoder, authHttp.LoginPost))
		router.Get("/logout", kithttp.HandlerWrap(resEncoder, authHttp.LoginOutGet))
		router.Post("/logout", kithttp.HandlerWrap(resEncoder, authHttp.Logout))
		router.Get("/consent", kithttp.HandlerWrap(resEncoder, authHttp.ConsentGet))
		router.Post("/consent", kithttp.HandlerWrap(resEncoder, authHttp.Consent))

		kithttp.MountBlob(srv, "", biz.UserAvatarPath, vfs)

		srv.HandlePrefix("/v1/auth/web", http.StripPrefix("/v1/auth/web", router))

		v12.RegisterUserAdminServiceHTTPServer(srv, user)
		v12.RegisterUserServiceHTTPServer(srv, user)

		v13.RegisterAccountHTTPServer(srv, account)

		route := srv.Route("/")

		route.POST("/v1/account/avatar", account.UpdateAvatar)
		route.POST("/v1/user/avatar", user.UpdateAvatar)

		v14.RegisterAuthHTTPServer(srv, auth)
		v14.RegisterWeChatAuthServiceHTTPServer(srv, weChatAuth)
		v1.RegisterRoleServiceHTTPServer(srv, role)
		v15.RegisterPermissionServiceHTTPServer(srv, permission)

		swaggerRouter := chi.NewRouter()
		swaggerRouter.Use(
			kithttp.MiddlewareConvert(errEncoder, middleware...))
		const apiPrefix = "/v1/user/dev/swagger"
		swaggerRouter.Handle(apiPrefix+"*", http.StripPrefix(apiPrefix, swaggerui.Handler(spec)))

		srv.HandlePrefix(apiPrefix, swaggerRouter)
	})
}

func NewGrpcServerRegister(
	user *UserService,
	userInternal *UserInternalService,
	account *AccountService,
	auth *AuthService,
	weChatAuth *WeChatAuthService,
	role *RoleService,
	permission *PermissionService) kitgrpc.ServiceRegister {
	return kitgrpc.ServiceRegisterFunc(func(srv *grpc.Server, middleware ...middleware.Middleware) {
		v12.RegisterUserAdminServiceServer(srv, user)
		v12.RegisterUserServiceServer(srv, user)
		v12.RegisterUserInternalServiceServer(srv, userInternal)
		v13.RegisterAccountServer(srv, account)
		v14.RegisterAuthServer(srv, auth)
		v14.RegisterWeChatAuthServiceServer(srv, weChatAuth)
		v1.RegisterRoleServiceServer(srv, role)
		v15.RegisterPermissionServiceServer(srv, permission)
		v15.RegisterPermissionInternalServiceServer(srv, permission)
	})
}
