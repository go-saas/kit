package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	sapi "github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/authn/jwt"
	"github.com/go-saas/kit/pkg/conf"
	"github.com/go-saas/kit/pkg/localize"
	"github.com/go-saas/kit/pkg/logging"
	"github.com/go-saas/kit/pkg/server"
	"github.com/go-saas/kit/pkg/server/common"
	kitgrpc "github.com/go-saas/kit/pkg/server/grpc"
	"github.com/go-saas/kit/saas/api"
	v12 "github.com/go-saas/kit/saas/api/plan/v1"
	v1 "github.com/go-saas/kit/saas/api/tenant/v1"
	uapi "github.com/go-saas/kit/user/api"
	"github.com/go-saas/saas"
	"github.com/go-saas/saas/http"
	uow2 "github.com/go-saas/uow"
	kuow "github.com/go-saas/uow/kratos"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Services, tokenizer jwt.Tokenizer, ts saas.TenantStore, uowMgr uow2.Manager,
	mOpt *http.WebMultiTenancyOption, apiOpt *sapi.Option,
	userTenant *uapi.UserTenantContrib,
	validator sapi.TrustedContextValidator,
	register []kitgrpc.ServiceRegister,
	logger log.Logger) *kitgrpc.Server {
	m := []middleware.Middleware{server.Recovery(),
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
		)),
	}
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			m...,
		),
	}
	cfg := common.GetConf(c, api.ServiceName)
	opts = kitgrpc.PatchOpts(logger, opts, cfg)
	srv := kitgrpc.NewServer(cfg, opts...)
	kitgrpc.ChainServiceRegister(register...).Register(srv.Server, m...)

	return srv
}
