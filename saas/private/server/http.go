package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
	api2 "github.com/goxiaoy/go-saas-kit/pkg/api"
	"github.com/goxiaoy/go-saas-kit/pkg/authn/jwt"
	"github.com/goxiaoy/go-saas-kit/pkg/blob"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	"github.com/goxiaoy/go-saas-kit/pkg/server"
	"github.com/goxiaoy/go-saas-kit/pkg/uow"
	"github.com/goxiaoy/go-saas-kit/saas/api"
	v1 "github.com/goxiaoy/go-saas-kit/saas/api/tenant/v1"
	conf2 "github.com/goxiaoy/go-saas-kit/saas/private/conf"
	"github.com/goxiaoy/go-saas-kit/saas/private/service"
	"github.com/goxiaoy/go-saas-kit/user/remote"
	"github.com/goxiaoy/go-saas/common"
	http2 "github.com/goxiaoy/go-saas/common/http"
	uow2 "github.com/goxiaoy/uow"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer(c *conf.Services,
	sCfg *conf.Security,
	tokenizer jwt.Tokenizer,
	ts common.TenantStore,
	uowMgr uow2.Manager,
	tenant *service.TenantService,
	mOpt *http2.WebMultiTenancyOption,
	apiOpt *api2.Option,
	reqDecoder http.DecodeRequestFunc,
	resEncoder http.EncodeResponseFunc,
	errEncoder http.EncodeErrorFunc,
	factory blob.Factory,
	dataCfg *conf2.Data,
	logger log.Logger,
	userTenant *remote.UserTenantContributor) *http.Server {
	var opts []http.ServerOption
	opts = server.PatchHttpOpts(logger, opts, api.ServiceName, c, sCfg, reqDecoder, resEncoder, errEncoder)
	opts = append(opts, []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			tracing.Server(),
			logging.Server(logger),
			metrics.Server(),
			validate.Validator(),
			jwt.ServerExtractAndAuth(tokenizer, logger),
			api2.ServerMiddleware(apiOpt, logger),
			server.Saas(mOpt, ts, func(o *common.TenantResolveOption) {
				o.AppendContributors(userTenant)
			}),
			uow.Uow(logger, uowMgr),
		),
	}...)
	srv := http.NewServer(opts...)

	route := srv.Route("/")

	route.POST("/v1/saas/tenant/logo", tenant.UpdateLogo)
	server.HandleBlobs("", dataCfg.Blobs, srv, factory)
	v1.RegisterTenantServiceHTTPServer(srv, tenant)
	return srv
}
