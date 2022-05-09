package server

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	sapi "github.com/goxiaoy/go-saas-kit/pkg/api"
	"github.com/goxiaoy/go-saas-kit/pkg/authn/jwt"
	"github.com/goxiaoy/go-saas-kit/pkg/authn/session"
	"github.com/goxiaoy/go-saas-kit/pkg/blob"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	localize "github.com/goxiaoy/go-saas-kit/pkg/i18n"
	"github.com/goxiaoy/go-saas-kit/pkg/server"
	"github.com/goxiaoy/go-saas-kit/pkg/uow"
	"github.com/goxiaoy/go-saas-kit/user/api"
	v13 "github.com/goxiaoy/go-saas-kit/user/api/account/v1"
	v14 "github.com/goxiaoy/go-saas-kit/user/api/auth/v1"
	v15 "github.com/goxiaoy/go-saas-kit/user/api/permission/v1"
	v1 "github.com/goxiaoy/go-saas-kit/user/api/role/v1"
	v12 "github.com/goxiaoy/go-saas-kit/user/api/user/v1"
	"github.com/goxiaoy/go-saas-kit/user/i18n"
	"github.com/goxiaoy/go-saas-kit/user/private/biz"
	conf2 "github.com/goxiaoy/go-saas-kit/user/private/conf"
	uhttp "github.com/goxiaoy/go-saas-kit/user/private/server/http"
	"github.com/goxiaoy/go-saas-kit/user/private/service"
	"github.com/goxiaoy/go-saas/common"
	shttp "github.com/goxiaoy/go-saas/common/http"
	uow2 "github.com/goxiaoy/uow"
	"net/http"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer(c *conf.Services,
	sCfg *conf.Security,
	tokenizer jwt.Tokenizer,
	uowMgr uow2.Manager,
	mOpt *shttp.WebMultiTenancyOption,
	apiOpt *sapi.Option,
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
	dataCfg *conf2.Data,
	factory blob.Factory,
	userTenant *service.UserTenantContributor,
	validator sapi.TrustedContextValidator,
	refreshProvider session.RefreshTokenProvider,
) *khttp.Server {
	var opts []khttp.ServerOption
	opts = server.PatchHttpOpts(logger, opts, api.ServiceName, c, sCfg, reqDecoder, resEncoder, errEncoder,
		session.Auth(sCfg, validator),
		session.Refresh(errEncoder, refreshProvider, validator),
	)
	middlewares := middleware.Chain(recovery.Recovery(),
		tracing.Server(),
		logging.Server(logger),
		metrics.Server(),
		validate.Validator(),
		localize.I18N(i18n.Files...),
		jwt.ServerExtractAndAuth(tokenizer, logger),
		sapi.ServerPropagation(apiOpt, validator, logger),
		server.Saas(mOpt, ts, validator, func(o *common.TenantResolveOption) {
			o.AppendContributors(userTenant)
		}),
		uow.Uow(logger, uowMgr))
	opts = append(opts, []khttp.ServerOption{
		khttp.Middleware(middlewares),
	}...)

	router := chi.NewRouter()

	//global filter
	router.Use(
		server.MiddlewareConvert(errEncoder, middlewares))

	router.Group(func(router chi.Router) {
		router.Get("/login", server.HandlerWrap(resEncoder, authHttp.LoginGet))
		router.Post("/login", server.HandlerWrap(resEncoder, authHttp.LoginPost))
		router.Get("/logout", server.HandlerWrap(resEncoder, authHttp.LoginOutGet))
		router.Post("/logout", server.HandlerWrap(resEncoder, authHttp.Logout))
		router.Get("/consent", server.HandlerWrap(resEncoder, authHttp.ConsentGet))
		router.Post("/consent", server.HandlerWrap(resEncoder, authHttp.Consent))
	})

	srv := khttp.NewServer(opts...)
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
	return srv
}

func NewRefreshTokenProvider(sign *biz.SignInManager) session.RefreshTokenProvider {
	return session.RefreshTokenProviderFunc(func(ctx context.Context, token, userId string) (err error) {
		return sign.RefreshSignIn(ctx, token)
	})
}
