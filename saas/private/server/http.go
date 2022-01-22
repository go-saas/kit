package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
	api2 "github.com/goxiaoy/go-saas-kit/pkg/api"
	"github.com/goxiaoy/go-saas-kit/pkg/authn/jwt"
	"github.com/goxiaoy/go-saas-kit/pkg/authn/middleware/authentication"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	"github.com/goxiaoy/go-saas-kit/pkg/server"
	"github.com/goxiaoy/go-saas-kit/pkg/uow"
	"github.com/goxiaoy/go-saas-kit/saas/api"
	v1 "github.com/goxiaoy/go-saas-kit/saas/api/tenant/v1"
	"github.com/goxiaoy/go-saas-kit/saas/private/service"
	"github.com/goxiaoy/go-saas/common"
	http2 "github.com/goxiaoy/go-saas/common/http"
	"github.com/goxiaoy/go-saas/kratos/saas"
	uow2 "github.com/goxiaoy/uow"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer(c *conf.Services, sCfg *conf.Security, tokenizer jwt.Tokenizer, ts common.TenantStore, uowMgr uow2.Manager, tenant *service.TenantService, mOpt *http2.WebMultiTenancyOption, apiOpt *api2.Option, logger log.Logger) *http.Server {
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
			selector.Server(uow.Uow(logger, uowMgr)).Match(uow.DefaultUseOperation).Build(),
		),
	}...)
	srv := http.NewServer(opts...)

	v1.RegisterTenantServiceHTTPServer(srv, tenant)
	return srv
}
