package server

import (
	_ "embed"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/mux"
	"github.com/goxiaoy/go-saas-kit/payment/api"
	v12 "github.com/goxiaoy/go-saas-kit/payment/api/post/v1"
	"github.com/goxiaoy/go-saas-kit/payment/private/conf"
	"github.com/goxiaoy/go-saas-kit/payment/private/service"
	api2 "github.com/goxiaoy/go-saas-kit/pkg/api"
	sapi "github.com/goxiaoy/go-saas-kit/pkg/api"
	"github.com/goxiaoy/go-saas-kit/pkg/authn/jwt"
	"github.com/goxiaoy/go-saas-kit/pkg/blob"
	conf2 "github.com/goxiaoy/go-saas-kit/pkg/conf"
	"github.com/goxiaoy/go-saas-kit/pkg/server"
	"github.com/goxiaoy/go-saas-kit/pkg/uow"
	"github.com/goxiaoy/go-saas/common"
	shttp "github.com/goxiaoy/go-saas/common/http"
	uow2 "github.com/goxiaoy/uow"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer(c *conf2.Services,
	sCfg *conf2.Security,
	tokenizer jwt.Tokenizer,
	ts common.TenantStore,
	uowMgr uow2.Manager,
	reqDecoder khttp.DecodeRequestFunc,
	resEncoder khttp.EncodeResponseFunc,
	errEncoder khttp.EncodeErrorFunc,
	factory blob.Factory,
	dataCfg *conf.Data,
	mOpt *shttp.WebMultiTenancyOption,
	apiOpt *api2.Option,
	post *service.PostServiceService,
	validator sapi.TrustedContextValidator,
	logger log.Logger) *khttp.Server {
	var opts []khttp.ServerOption
	opts = server.PatchHttpOpts(logger, opts, api.ServiceName, c, sCfg, reqDecoder, resEncoder, errEncoder)
	opts = append(opts, []khttp.ServerOption{
		khttp.Middleware(
			recovery.Recovery(),
			tracing.Server(),
			logging.Server(logger),
			metrics.Server(),
			validate.Validator(),
			jwt.ServerExtractAndAuth(tokenizer, logger),
			sapi.ServerPropagation(apiOpt, validator, logger),
			server.Saas(mOpt, ts, validator),
			uow.Uow(logger, uowMgr),
		),
	}...)
	srv := khttp.NewServer(opts...)

	//handle swagger
	router := mux.NewRouter()
	router.Use(server.MiddlewareConvert(errEncoder,
		recovery.Recovery(),
		tracing.Server(),
		logging.Server(logger),
		metrics.Server()))

	server.HandleBlobs("", dataCfg.Blobs, srv, factory)

	v12.RegisterPostServiceHTTPServer(srv, post)
	return srv
}
