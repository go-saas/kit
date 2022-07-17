package service

import (
	_ "embed"
	"github.com/flowchartsman/swaggerui"
	"github.com/go-chi/chi/v5"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-saas/kit/pkg/authz/authz"
	"github.com/go-saas/kit/pkg/blob"
	kconf "github.com/go-saas/kit/pkg/conf"
	kitdi "github.com/go-saas/kit/pkg/di"
	"github.com/go-saas/kit/pkg/server"
	"github.com/go-saas/kit/user/api"
	v13 "github.com/go-saas/kit/user/api/account/v1"
	v14 "github.com/go-saas/kit/user/api/auth/v1"
	v15 "github.com/go-saas/kit/user/api/permission/v1"
	v1 "github.com/go-saas/kit/user/api/role/v1"
	v12 "github.com/go-saas/kit/user/api/user/v1"
	uhttp "github.com/go-saas/kit/user/private/service/http"
	"github.com/goava/di"
	client "github.com/ory/hydra-client-go"
	"net/http"
)

//go:embed openapi/api.swagger.json
var spec []byte

// ProviderSet is service providers.
var ProviderSet = kitdi.NewSet(
	NewGrpcServerRegister,
	NewHttpServerRegister,
	NewUserRoleContrib,
	kitdi.NewProvider(NewUserService, di.As(new(v12.UserServiceServer))),

	kitdi.NewProvider(NewAccountService, di.As(new(v13.AccountServer))),

	kitdi.NewProvider(NewAuthService, di.As(new(v14.AuthServer))),

	kitdi.NewProvider(NewRoleServiceService, di.As(new(v1.RoleServiceServer))),
	kitdi.NewProvider(NewPermissionService, di.As(new(v15.PermissionServiceServer))),
	NewHydra,
	api.NewUserTenantContrib,
	api.NewRefreshProvider,
	uhttp.NewAuth)

func NewHttpServerRegister(user *UserService,
	resEncoder khttp.EncodeResponseFunc,
	errEncoder khttp.EncodeErrorFunc,
	account *AccountService,
	auth *AuthService,
	role *RoleService,
	permission *PermissionService,
	authHttp *uhttp.Auth,
	dataCfg *kconf.Data,
	authzSrv authz.Service,
	factory blob.Factory) server.HttpServiceRegister {
	return server.HttpServiceRegisterFunc(func(srv *khttp.Server, middleware ...middleware.Middleware) {

		router := chi.NewRouter()

		//global filter
		router.Use(
			server.MiddlewareConvert(errEncoder, middleware...))

		router.Get("/login", server.HandlerWrap(resEncoder, authHttp.LoginGet))
		router.Post("/login", server.HandlerWrap(resEncoder, authHttp.LoginPost))
		router.Get("/logout", server.HandlerWrap(resEncoder, authHttp.LoginOutGet))
		router.Post("/logout", server.HandlerWrap(resEncoder, authHttp.Logout))
		router.Get("/consent", server.HandlerWrap(resEncoder, authHttp.ConsentGet))
		router.Post("/consent", server.HandlerWrap(resEncoder, authHttp.Consent))

		server.HandleBlobs("", dataCfg.Blobs, srv, factory)
		srv.HandlePrefix("/v1/auth/web", http.StripPrefix("/v1/auth/web", router))

		v12.RegisterUserServiceHTTPServer(srv, user)

		v13.RegisterAccountHTTPServer(srv, account)

		route := srv.Route("/")

		route.POST("/v1/account/avatar", account.UpdateAvatar)
		route.POST("/v1/user/avatar", user.UpdateAvatar)

		v14.RegisterAuthHTTPServer(srv, auth)
		v1.RegisterRoleServiceHTTPServer(srv, role)
		v15.RegisterPermissionServiceHTTPServer(srv, permission)

		swaggerRouter := chi.NewRouter()
		swaggerRouter.Use(
			server.MiddlewareConvert(errEncoder, middleware...))
		const apiPrefix = "/v1/user/dev/swagger"
		swaggerRouter.Handle(apiPrefix+"*", http.StripPrefix(apiPrefix, server.AuthzGuardian(
			authzSrv, authz.RequirementList{
				authz.NewRequirement(authz.NewEntityResource("dev", "user"), authz.AnyAction),
			}, errEncoder, swaggerui.Handler(spec),
		)))

		srv.HandlePrefix(apiPrefix, swaggerRouter)
	})
}

func NewGrpcServerRegister(user *UserService,
	account *AccountService,
	auth *AuthService,
	role *RoleService,
	permission *PermissionService) server.GrpcServiceRegister {
	return server.GrpcServiceRegisterFunc(func(srv *grpc.Server, middleware ...middleware.Middleware) {
		v12.RegisterUserServiceServer(srv, user)
		v13.RegisterAccountServer(srv, account)
		v14.RegisterAuthServer(srv, auth)
		v1.RegisterRoleServiceServer(srv, role)
		v15.RegisterPermissionServiceServer(srv, permission)
	})
}

func NewHydra(c *kconf.Security) *client.APIClient {
	cfg := client.NewConfiguration()
	cfg.Servers = client.ServerConfigurations{
		{
			URL: c.Oidc.Hydra.AdminUrl,
		},
	}
	return client.NewAPIClient(cfg)
}
