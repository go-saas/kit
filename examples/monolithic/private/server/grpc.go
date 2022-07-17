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
	uapi "github.com/go-saas/kit/user/api"
	"github.com/go-saas/saas"
	shttp "github.com/go-saas/saas/http"
	uow2 "github.com/go-saas/uow"
	kuow "github.com/go-saas/uow/kratos"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(
	c *conf.Services,
	tokenizer jwt.Tokenizer,
	ts saas.TenantStore,
	uowMgr uow2.Manager,
	mOpt *shttp.WebMultiTenancyOption,
	apiOpt *sapi.Option,
	logger log.Logger,
	validator sapi.TrustedContextValidator,
	userTenant *uapi.UserTenantContrib,
	register []server.GrpcServiceRegister,
) *grpc.Server {
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
		kuow.Uow(uowMgr, kuow.WithLogger(logger)),
	}
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			m...,
		),
	}
	opts = server.PatchGrpcOpts(logger, opts, uapi.ServiceName, c)
	srv := grpc.NewServer(opts...)

	server.ChainGrpcServiceRegister(register...).Register(srv, middleware.Chain(m...))

	return srv
}
