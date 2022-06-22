package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	sapi "github.com/goxiaoy/go-saas-kit/pkg/api"
	"github.com/goxiaoy/go-saas-kit/pkg/authn/jwt"
	conf2 "github.com/goxiaoy/go-saas-kit/pkg/conf"
	"github.com/goxiaoy/go-saas-kit/pkg/localize"
	"github.com/goxiaoy/go-saas-kit/pkg/logging"
	"github.com/goxiaoy/go-saas-kit/pkg/server"
	"github.com/goxiaoy/go-saas-kit/pkg/uow"
	"github.com/goxiaoy/go-saas-kit/sys/api"
	"github.com/goxiaoy/go-saas-kit/sys/i18n"
	"github.com/goxiaoy/go-saas-kit/sys/private/service"
	uow2 "github.com/goxiaoy/uow"
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
		uow.Uow(logger, uowMgr),
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
