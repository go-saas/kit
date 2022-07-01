package api

import (
	"github.com/dtm-labs/dtmgrpc"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/conf"
)

func NewInit(client *conf.Client, opt *api.Option, tokenMgr api.TokenManager, logger klog.Logger) {
	m := api.ClientPropagation(client, opt, tokenMgr, logger)
	dtmgrpc.AddUnaryInterceptor(api.UnaryClientInterceptor([]middleware.Middleware{m}, 0, nil))
}
