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
	"github.com/goxiaoy/go-saas-kit/pkg/uow"
	v13 "github.com/goxiaoy/go-saas-kit/user/api/account/v1"
	v14 "github.com/goxiaoy/go-saas-kit/user/api/auth/v1"
	v12 "github.com/goxiaoy/go-saas-kit/user/api/user/v1"
	"github.com/goxiaoy/go-saas-kit/user/internal_/conf"
	"github.com/goxiaoy/go-saas-kit/user/internal_/service"
	"github.com/goxiaoy/go-saas/common"
	"github.com/goxiaoy/go-saas/kratos/middleware"
	uow2 "github.com/goxiaoy/uow"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer(c *conf.Server, tokenizer jwt.Tokenizer, uowMgr uow2.Manager, ts common.TenantStore, user *service.UserService, account *service.AccountService, auth *service.AuthService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			tracing.Server(),
			logging.Server(logger),
			metrics.Server(),
			validate.Validator(),
			authentication.ServerExtractAndAuth(logger, tokenizer),
			uow.Uow(logger, uowMgr),
			middleware.MultiTenancy(nil, nil, ts),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	openAPIhandler := openapiv2.NewHandler()
	srv := http.NewServer(opts...)
	srv.HandlePrefix("/q/", openAPIhandler)
	v12.RegisterUserServiceHTTPServer(srv, user)
	v13.RegisterAccountHTTPServer(srv, account)
	v14.RegisterAuthHTTPServer(srv, auth)
	return srv
}
