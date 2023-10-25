package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
	sapi "github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/authn/jwt"
	"github.com/go-saas/kit/pkg/conf"
	"github.com/go-saas/kit/pkg/localize"
	"github.com/go-saas/kit/pkg/logging"
	"github.com/go-saas/kit/pkg/server"
	"github.com/go-saas/kit/pkg/server/common"
	kithttp "github.com/go-saas/kit/pkg/server/http"
	"github.com/go-saas/kit/saas/api"
	v12 "github.com/go-saas/kit/saas/api/plan/v1"
	v1 "github.com/go-saas/kit/saas/api/tenant/v1"
	uapi "github.com/go-saas/kit/user/api"
	"github.com/go-saas/saas"
	http2 "github.com/go-saas/saas/http"
	uow2 "github.com/go-saas/uow"
	kuow "github.com/go-saas/uow/kratos"
)

// NewHTTPServer new a HTTP kithttp.
func NewHTTPServer(c *conf.Services,
	sCfg *conf.Security,
	tokenizer jwt.Tokenizer,
	ts saas.TenantStore,
	uowMgr uow2.Manager,
	mOpt *http2.WebMultiTenancyOption,
	apiOpt *sapi.Option,
	reqDecoder http.DecodeRequestFunc,
	resEncoder http.EncodeResponseFunc,
	errEncoder http.EncodeErrorFunc,
	logger log.Logger,
	validator sapi.TrustedContextValidator,
	userTenant *uapi.UserTenantContrib,
	register []kithttp.ServiceRegister) *kithttp.Server {
	var opts []http.ServerOption
	cfg := common.GetConf(c, api.ServiceName)
	opts = kithttp.PatchOpts(logger, opts, cfg, sCfg, reqDecoder, resEncoder, errEncoder)
	m := []middleware.Middleware{
		server.Recovery(),
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
		kuow.Uow(uowMgr, kuow.WithForceSkipOp(
			v1.TenantService_CreateTenant_FullMethodName,
			v12.PlanService_CreatePlan_FullMethodName,
			v12.PlanService_UpdatePlan_FullMethodName,
			v12.PlanService_DeletePlan_FullMethodName,
		))}

	opts = append(opts, []http.ServerOption{
		http.Middleware(
			m...,
		),
	}...)
	srv := kithttp.NewServer(cfg, opts...)
	kithttp.ChainServiceRegister(register...).Register(srv.Server, m...)

	return srv
}
