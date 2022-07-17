package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	sapi "github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/authn/jwt"
	"github.com/go-saas/kit/pkg/authn/session"
	"github.com/go-saas/kit/pkg/conf"
	"github.com/go-saas/kit/pkg/localize"
	"github.com/go-saas/kit/pkg/logging"
	"github.com/go-saas/kit/pkg/server"
	"github.com/go-saas/kit/user/api"
	"github.com/go-saas/saas"
	shttp "github.com/go-saas/saas/http"
	uow2 "github.com/go-saas/uow"
	kuow "github.com/go-saas/uow/kratos"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer(c *conf.Services,
	sCfg *conf.Security,
	tokenizer jwt.Tokenizer,
	uowMgr uow2.Manager,
	mOpt *shttp.WebMultiTenancyOption,
	apiOpt *sapi.Option,
	ts saas.TenantStore,
	reqDecoder khttp.DecodeRequestFunc,
	resEncoder khttp.EncodeResponseFunc,
	errEncoder khttp.EncodeErrorFunc,
	logger log.Logger,
	userTenant *api.UserTenantContrib,
	validator sapi.TrustedContextValidator,
	refreshProvider session.RefreshTokenProvider,
	register []server.HttpServiceRegister,
) *khttp.Server {
	var opts []khttp.ServerOption
	opts = server.PatchHttpOpts(logger, opts, api.ServiceName, c, sCfg, reqDecoder, resEncoder, errEncoder,
		session.Auth(sCfg, validator),
		session.Refresh(errEncoder, refreshProvider, validator),
	)
	middlewares := []middleware.Middleware{server.Recovery(),
		tracing.Server(),
		logging.Server(logger),
		metrics.Server(),
		validate.Validator(),
		localize.I18N(),
		jwt.ServerExtractAndAuth(tokenizer, logger),
		sapi.ServerPropagation(apiOpt, validator, logger),
		server.Saas(mOpt, ts, validator, func(o *saas.TenantResolveOption) {
			o.AppendContribs(userTenant)
		}),
		kuow.Uow(uowMgr, kuow.WithLogger(logger)),
	}
	opts = append(opts, []khttp.ServerOption{
		khttp.Middleware(middlewares...),
	}...)

	srv := khttp.NewServer(opts...)

	server.ChainHttpServiceRegister(register...).Register(srv, middlewares...)
	return srv
}
