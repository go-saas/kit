package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	sapi "github.com/goxiaoy/go-saas-kit/pkg/api"
	"github.com/goxiaoy/go-saas-kit/pkg/authn/jwt"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	"github.com/goxiaoy/go-saas-kit/pkg/localize"
	"github.com/goxiaoy/go-saas-kit/pkg/server"
	"github.com/goxiaoy/go-saas-kit/pkg/uow"
	"github.com/goxiaoy/go-saas-kit/user/api"
	"github.com/goxiaoy/go-saas-kit/user/i18n"
	"github.com/goxiaoy/go-saas-kit/user/private/service"
	"github.com/goxiaoy/go-saas/common"
	http2 "github.com/goxiaoy/go-saas/common/http"
	uow2 "github.com/goxiaoy/uow"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(
	c *conf.Services,
	tokenizer jwt.Tokenizer,
	ts common.TenantStore,
	uowMgr uow2.Manager,
	mOpt *http2.WebMultiTenancyOption,
	apiOpt *sapi.Option,
	logger log.Logger,
	validator sapi.TrustedContextValidator,
	userTenant *service.UserTenantContributor,
	register GrpcServerRegister,
) *grpc.Server {
	m := middleware.Chain(
		server.Recovery(),
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
