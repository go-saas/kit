package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/sessions"
	api2 "github.com/goxiaoy/go-saas-kit/pkg/api"
	"github.com/goxiaoy/go-saas-kit/pkg/authn/jwt"
	"github.com/goxiaoy/go-saas-kit/pkg/authn/middleware/authentication"
	"github.com/goxiaoy/go-saas-kit/pkg/authn/session"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	"github.com/goxiaoy/go-saas-kit/pkg/server"
	"github.com/goxiaoy/go-saas-kit/pkg/uow"
	"github.com/goxiaoy/go-saas-kit/user/api"
	v13 "github.com/goxiaoy/go-saas-kit/user/api/account/v1"
	v14 "github.com/goxiaoy/go-saas-kit/user/api/auth/v1"
	v15 "github.com/goxiaoy/go-saas-kit/user/api/permission/v1"
	v1 "github.com/goxiaoy/go-saas-kit/user/api/role/v1"
	v12 "github.com/goxiaoy/go-saas-kit/user/api/user/v1"
	uhttp "github.com/goxiaoy/go-saas-kit/user/private/server/http"
	"github.com/goxiaoy/go-saas-kit/user/private/service"
	"github.com/goxiaoy/go-saas/common"
	shttp "github.com/goxiaoy/go-saas/common/http"
	"github.com/goxiaoy/go-saas/kratos/saas"
	uow2 "github.com/goxiaoy/uow"
	"net/http"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer(c *conf.Services,
	sCfg *conf.Security,
	tokenizer jwt.Tokenizer,
	uowMgr uow2.Manager,
	mOpt *shttp.WebMultiTenancyOption,
	apiOpt *api2.Option,
	ts common.TenantStore,
	reqDecoder khttp.DecodeRequestFunc,
	resEncoder khttp.EncodeResponseFunc,
	errEncoder khttp.EncodeErrorFunc,
	logger log.Logger,
	user *service.UserService,
	account *service.AccountService,
	auth *service.AuthService,
	role *service.RoleService,
	permission *service.PermissionService,
	authHttp *uhttp.Auth,
	errorHandler server.ErrorHandler,
	sessionStore sessions.Store,
) *khttp.Server {
	var opts []khttp.ServerOption
	opts = server.PatchHttpOpts(logger, opts, api.ServiceName, c, sCfg, reqDecoder, resEncoder, errEncoder,
		//extract from session cookie
		session.Auth(sessionStore, sCfg))

	opts = append(opts, []khttp.ServerOption{
		khttp.Middleware(
			recovery.Recovery(),
			tracing.Server(),
			logging.Server(logger),
			metrics.Server(),
			validate.Validator(),
			authentication.ServerExtractAndAuth(tokenizer, logger),
			saas.Server(mOpt, nil, ts),
			api2.ServerMiddleware(apiOpt),
			uow.Uow(logger, uowMgr),
		),
	}...)

	router := chi.NewRouter()

	//global filter
	router.Use(
		server.MiddlewareConvert(recovery.Recovery(),
			tracing.Server(),
			logging.Server(logger),
			metrics.Server(),
			validate.Validator(),
			authentication.ServerExtractAndAuth(tokenizer, logger)),

		server.MiddlewareConvert(
			saas.Server(mOpt, nil, ts),
			api2.ServerMiddleware(apiOpt),
			uow.Uow(logger, uowMgr),
		))

	router.Group(func(router chi.Router) {
		router.Get("/login", errorHandler.Wrap(authHttp.LoginGet).ServeHTTP)
		router.Post("/login", errorHandler.Wrap(authHttp.LoginPost).ServeHTTP)
	})

	srv := khttp.NewServer(opts...)

	srv.HandlePrefix("/v1/auth/web", http.StripPrefix("/v1/auth/web", router))

	v12.RegisterUserServiceHTTPServer(srv, user)
	v13.RegisterAccountHTTPServer(srv, account)
	v14.RegisterAuthHTTPServer(srv, auth)
	v1.RegisterRoleServiceHTTPServer(srv, role)
	v15.RegisterPermissionServiceHTTPServer(srv, permission)
	return srv
}
