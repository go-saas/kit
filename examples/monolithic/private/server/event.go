package server

import (
	klog "github.com/go-kratos/kratos/v2/log"
	kitconf "github.com/goxiaoy/go-saas-kit/pkg/conf"
	"github.com/goxiaoy/go-saas-kit/pkg/dal"
	"github.com/goxiaoy/go-saas-kit/pkg/event"
	"github.com/goxiaoy/go-saas-kit/pkg/event/trace"
	saasbiz "github.com/goxiaoy/go-saas-kit/saas/private/biz"
	ubiz "github.com/goxiaoy/go-saas-kit/user/private/biz"
	uow2 "github.com/goxiaoy/uow"
)

func NewEventServer(c *kitconf.Data, conn dal.ConnName, logger klog.Logger, uowMgr uow2.Manager, tenantReady saasbiz.TenantReadyEventHandler, tenantSeed ubiz.TenantSeedEventHandler) *event.ConsumerFactoryServer {
	e := c.Endpoints.GetEventMergedDefault(string(conn))
	srv := event.NewConsumerFactoryServer(e)
	srv.Use(event.ConsumerRecover(event.WithLogger(logger)), trace.Receive(), event.Logging(logger), event.ConsumerUow(uowMgr))
	srv.Append(tenantReady)
	srv.Append(tenantSeed)
	return srv
}
