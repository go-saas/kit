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
	conf2 "github.com/go-saas/kit/pkg/conf"
	"github.com/go-saas/kit/pkg/localize"
	"github.com/go-saas/kit/pkg/logging"
	"github.com/go-saas/kit/pkg/server"
	"github.com/go-saas/kit/sys/api"
	"github.com/go-saas/kit/sys/i18n"
	"github.com/go-saas/kit/sys/private/service"
	uow2 "github.com/go-saas/uow"
	kuow "github.com/go-saas/uow/kratos"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(
	c *conf2.Services,
	tokenizer jwt.Tokenizer,
	uowMgr uow2.Manager,
	apiOpt *sapi.Option,
	logger log.Logger,
	validator sapi.TrustedContextValidator,
	register service.GrpcServerRegister,
) *grpc.Server {
	m := middleware.Chain(
		server.Recovery(),
		tracing.Server(),
		logging.Server(logger),
		metrics.Server(),
		validate.Validator(),
		localize.I18N(i18n.Files...),
		jwt.ServerExtractAndAuth(tokenizer, logger),
		sapi.ServerPropagation(apiOpt, validator, logger),
		kuow.Uow(uowMgr, kuow.WithLogger(logger)),
	)
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			m,
		),
	}
	opts = server.PatchGrpcOpts(logger, opts, api.ServiceName, c)
	srv := grpc.NewServer(opts...)
	register.Register(srv, m)

	return srv
}
