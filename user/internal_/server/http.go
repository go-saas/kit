package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/swagger-api/openapiv2"
	"github.com/goxiaoy/go-saas-kit/auth/jwt"
	"github.com/goxiaoy/go-saas-kit/auth/middleware/authentication"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	"github.com/goxiaoy/go-saas-kit/pkg/server"
	"github.com/goxiaoy/go-saas-kit/pkg/uow"
	"github.com/goxiaoy/go-saas-kit/user/api"
	v13 "github.com/goxiaoy/go-saas-kit/user/api/account/v1"
	v14 "github.com/goxiaoy/go-saas-kit/user/api/auth/v1"
	v12 "github.com/goxiaoy/go-saas-kit/user/api/user/v1"
	"github.com/goxiaoy/go-saas-kit/user/internal_/service"
	"github.com/goxiaoy/go-saas/common"
	http2 "github.com/goxiaoy/go-saas/common/http"
	"github.com/goxiaoy/go-saas/kratos/middleware"
	uow2 "github.com/goxiaoy/uow"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer(c *conf.Services, tokenizer jwt.Tokenizer, uowMgr uow2.Manager, mOpt *http2.WebMultiTenancyOption, ts common.TenantStore, user *service.UserService, account *service.AccountService, auth *service.AuthService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			tracing.Server(),
			logging.Server(logger),
			metrics.Server(),
			validate.Validator(),
			authentication.ServerExtractAndAuth(tokenizer),
			uow.Uow(logger, uowMgr),
			middleware.Server(mOpt, nil, ts),
		),
	}
	opts = server.PatchHttpOpts(opts, api.ServiceName, c)

	openAPIhandler := openapiv2.NewHandler()
	srv := http.NewServer(opts...)
	srv.HandlePrefix("/q/", openAPIhandler)
	v12.RegisterUserServiceHTTPServer(srv, user)
	v13.RegisterAccountHTTPServer(srv, account)
	v14.RegisterAuthHTTPServer(srv, auth)
	return srv
}
