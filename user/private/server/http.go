package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
	api2 "github.com/goxiaoy/go-saas-kit/pkg/api"
	"github.com/goxiaoy/go-saas-kit/pkg/authn/jwt"
	authboss2 "github.com/goxiaoy/go-saas-kit/pkg/authn/middleware/authboss"
	"github.com/goxiaoy/go-saas-kit/pkg/authn/middleware/authentication"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	"github.com/goxiaoy/go-saas-kit/pkg/server"
	"github.com/goxiaoy/go-saas-kit/pkg/uow"
	"github.com/goxiaoy/go-saas-kit/user/api"
	v13 "github.com/goxiaoy/go-saas-kit/user/api/account/v1"
	v14 "github.com/goxiaoy/go-saas-kit/user/api/auth/v1"
	v15 "github.com/goxiaoy/go-saas-kit/user/api/permission/v1"
	v1 "github.com/goxiaoy/go-saas-kit/user/api/role/v1"
	v12 "github.com/goxiaoy/go-saas-kit/user/api/user/v1"
	"github.com/goxiaoy/go-saas-kit/user/private/service"
	"github.com/goxiaoy/go-saas/common"
	http2 "github.com/goxiaoy/go-saas/common/http"
	"github.com/goxiaoy/go-saas/kratos/saas"
	uow2 "github.com/goxiaoy/uow"
	"github.com/volatiletech/authboss/v3"
	http3 "net/http"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer(c *conf.Services,
	sCfg *conf.Security,
	tokenizer jwt.Tokenizer,
	uowMgr uow2.Manager,
	mOpt *http2.WebMultiTenancyOption,
	apiOpt *api2.Option,
	ts common.TenantStore,
	logger log.Logger,
	ab *authboss.Authboss,
	user *service.UserService,
	account *service.AccountService,
	auth *service.AuthService,
	role *service.RoleService,
	permission *service.PermissionService) *http.Server {
	var opts []http.ServerOption
	opts = server.PatchHttpOpts(logger, opts, api.ServiceName, c, sCfg)

	opts = append(opts, []http.ServerOption{
		http.Middleware(
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
			metrics.Server(), validate.Validator(),
			authentication.ServerExtractAndAuth(tokenizer, logger)),
		authboss2.PathFilter(ab),
		server.MiddlewareConvert(
			saas.Server(mOpt, nil, ts),
			api2.ServerMiddleware(apiOpt),
			uow.Uow(logger, uowMgr),
		),
		authboss.ModuleListMiddleware(ab))

	router.Group(func(router chi.Router) {
		router.Use(authboss.ModuleListMiddleware(ab))
		router.Mount("/", http3.StripPrefix("/auth", ab.Config.Core.Router))
	})

	srv := http.NewServer(opts...)

	srv.HandlePrefix("/auth", router)

	v12.RegisterUserServiceHTTPServer(srv, user)
	v13.RegisterAccountHTTPServer(srv, account)
	v14.RegisterAuthHTTPServer(srv, auth)
	v1.RegisterRoleServiceHTTPServer(srv, role)
	v15.RegisterPermissionServiceHTTPServer(srv, permission)
	return srv
}
