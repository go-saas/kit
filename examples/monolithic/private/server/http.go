package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	sapi "github.com/goxiaoy/go-saas-kit/pkg/api"
	"github.com/goxiaoy/go-saas-kit/pkg/authn/jwt"
	"github.com/goxiaoy/go-saas-kit/pkg/authn/session"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	"github.com/goxiaoy/go-saas-kit/pkg/localize"
	"github.com/goxiaoy/go-saas-kit/pkg/server"
	"github.com/goxiaoy/go-saas-kit/pkg/uow"
	"github.com/goxiaoy/go-saas-kit/user/api"
	"github.com/goxiaoy/go-saas-kit/user/i18n"
	"github.com/goxiaoy/go-saas-kit/user/private/service"
	"github.com/goxiaoy/go-saas/common"
	shttp "github.com/goxiaoy/go-saas/common/http"
	uow2 "github.com/goxiaoy/uow"
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
	userTenant *service.UserTenantContributor,
	validator sapi.TrustedContextValidator,
	refreshProvider session.RefreshTokenProvider,
	register HttpServerRegister,
) *khttp.Server {
	var opts []khttp.ServerOption
	opts = server.PatchHttpOpts(logger, opts, api.ServiceName, c, sCfg, reqDecoder, resEncoder, errEncoder,
		session.Auth(sCfg, validator),
		session.Refresh(errEncoder, refreshProvider, validator),
	)
	middlewares := middleware.Chain(server.Recovery(),
		tracing.Server(),
		logging.Server(logger),
		metrics.Server(),
		validate.Validator(),
		//TODO combine i18n
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

	srv := khttp.NewServer(opts...)

	register.Register(srv, middlewares)
	return srv
}
