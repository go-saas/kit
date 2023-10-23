package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	api2 "github.com/go-saas/kit/pkg/api"
	sapi "github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/authn/jwt"
	conf2 "github.com/go-saas/kit/pkg/conf"
	"github.com/go-saas/kit/pkg/localize"
	"github.com/go-saas/kit/pkg/server"
	"github.com/go-saas/kit/pkg/server/common"
	kitgrpc "github.com/go-saas/kit/pkg/server/grpc"
	"github.com/go-saas/kit/product/api"
	"github.com/go-saas/saas"
	shttp "github.com/go-saas/saas/http"
	uow2 "github.com/go-saas/uow"
	kuow "github.com/go-saas/uow/kratos"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(
	c *conf2.Services,
	tokenizer jwt.Tokenizer,
	ts saas.TenantStore,
	uowMgr uow2.Manager,
	mOpt *shttp.WebMultiTenancyOption,
	apiOpt *api2.Option,
	validator sapi.TrustedContextValidator,
	logger log.Logger,
	register []kitgrpc.ServiceRegister,
) *grpc.Server {
	m := []middleware.Middleware{
		recovery.Recovery(),
		tracing.Server(),
		logging.Server(logger),
		metrics.Server(),
		localize.I18N(),
		validate.Validator(),
		jwt.ServerExtractAndAuth(tokenizer, logger),
		sapi.ServerPropagation(apiOpt, validator, logger),
		server.Saas(mOpt, ts, validator),
		kuow.Uow(uowMgr),
	}
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			m...,
		),
	}
	cfg := common.GetConf(c, api.ServiceName)
	opts = kitgrpc.PatchOpts(logger, opts, cfg)
	srv := grpc.NewServer(opts...)
	kitgrpc.ChainServiceRegister(register...).Register(srv, m...)
	return srv
}
