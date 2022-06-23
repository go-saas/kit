package service

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/pkg/blob"
	kconf "github.com/goxiaoy/go-saas-kit/pkg/conf"
	"github.com/goxiaoy/go-saas-kit/pkg/server"
	"github.com/goxiaoy/go-saas-kit/user/api"
	v13 "github.com/goxiaoy/go-saas-kit/user/api/account/v1"
	v14 "github.com/goxiaoy/go-saas-kit/user/api/auth/v1"
	v15 "github.com/goxiaoy/go-saas-kit/user/api/permission/v1"
	v1 "github.com/goxiaoy/go-saas-kit/user/api/role/v1"
	v12 "github.com/goxiaoy/go-saas-kit/user/api/user/v1"
	uhttp "github.com/goxiaoy/go-saas-kit/user/private/service/http"
	client "github.com/ory/hydra-client-go"
	"net/http"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	NewGrpcServerRegister,
	NewHttpServerRegister,
	NewUserRoleContrib,
	NewUserService,
	wire.Bind(new(v12.UserServiceServer), new(*UserService)),
	NewAccountService,
	wire.Bind(new(v13.AccountServer), new(*AccountService)),
	NewAuthService,
	wire.Bind(new(v14.AuthServer), new(*AuthService)),
	NewRoleServiceService,
	wire.Bind(new(v1.RoleServiceServer), new(*RoleService)),
	NewPermissionService,
	wire.Bind(new(v15.PermissionServiceServer), new(*PermissionService)),
	NewHydra,
	api.NewUserTenantContrib,
	api.NewRefreshProvider,
	uhttp.NewAuth)

type HttpServerRegister server.HttpServiceRegister
type GrpcServerRegister server.GrpcServiceRegister

func NewHttpServerRegister(user *UserService,
	resEncoder khttp.EncodeResponseFunc,
	errEncoder khttp.EncodeErrorFunc,
	account *AccountService,
	auth *AuthService,
	role *RoleService,
	permission *PermissionService,
	authHttp *uhttp.Auth,
	dataCfg *kconf.Data,
	factory blob.Factory) HttpServerRegister {
	return server.HttpServiceRegisterFunc(func(srv *khttp.Server, middleware middleware.Middleware) {

		router := chi.NewRouter()

		//global filter
		router.Use(
			server.MiddlewareConvert(errEncoder, middleware))

		router.Group(func(router chi.Router) {
			router.Get("/login", server.HandlerWrap(resEncoder, authHttp.LoginGet))
			router.Post("/login", server.HandlerWrap(resEncoder, authHttp.LoginPost))
			router.Get("/logout", server.HandlerWrap(resEncoder, authHttp.LoginOutGet))
			router.Post("/logout", server.HandlerWrap(resEncoder, authHttp.Logout))
			router.Get("/consent", server.HandlerWrap(resEncoder, authHttp.ConsentGet))
			router.Post("/consent", server.HandlerWrap(resEncoder, authHttp.Consent))
		})

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
	})
}

func NewGrpcServerRegister(user *UserService,
	account *AccountService,
	auth *AuthService,
	role *RoleService,
	permission *PermissionService) GrpcServerRegister {
	return server.GrpcServiceRegisterFunc(func(srv *grpc.Server, middleware middleware.Middleware) {
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
