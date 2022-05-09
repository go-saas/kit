package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
	sapi "github.com/goxiaoy/go-saas-kit/pkg/api"
	"github.com/goxiaoy/go-saas-kit/pkg/authn/jwt"
	"github.com/goxiaoy/go-saas-kit/pkg/blob"
	conf2 "github.com/goxiaoy/go-saas-kit/pkg/conf"
	"github.com/goxiaoy/go-saas-kit/pkg/localize"
	"github.com/goxiaoy/go-saas-kit/pkg/server"
	"github.com/goxiaoy/go-saas-kit/pkg/uow"
	"github.com/goxiaoy/go-saas-kit/sys/api"
	v1 "github.com/goxiaoy/go-saas-kit/sys/api/menu/v1"
	"github.com/goxiaoy/go-saas-kit/sys/i18n"
	"github.com/goxiaoy/go-saas-kit/sys/private/conf"
	"github.com/goxiaoy/go-saas-kit/sys/private/service"
	"github.com/goxiaoy/go-saas/common"
	uow2 "github.com/goxiaoy/uow"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer(c *conf2.Services,
	sCfg *conf2.Security,
	tokenizer jwt.Tokenizer,
	ts common.TenantStore,
	uowMgr uow2.Manager,
	reqDecoder http.DecodeRequestFunc,
	resEncoder http.EncodeResponseFunc,
	errEncoder http.EncodeErrorFunc,
	factory blob.Factory,
	dataCfg *conf.Data,
	apiOpt *sapi.Option,
	logger log.Logger,
	validator sapi.TrustedContextValidator,
	menu *service.MenuService) *http.Server {
	var opts []http.ServerOption
	opts = server.PatchHttpOpts(logger, opts, api.ServiceName, c, sCfg, reqDecoder, resEncoder, errEncoder)
	opts = append(opts, []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			tracing.Server(),
			logging.Server(logger),
			metrics.Server(),
			validate.Validator(),
			localize.I18N(i18n.Files...),
			jwt.ServerExtractAndAuth(tokenizer, logger),
			sapi.ServerPropagation(apiOpt, validator, logger),
			uow.Uow(logger, uowMgr),
		),
	}...)
	srv := http.NewServer(opts...)

	server.HandleBlobs("", dataCfg.Blobs, srv, factory)

	v1.RegisterMenuServiceHTTPServer(srv, menu)

	return srv
}
