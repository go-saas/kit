package server

import (
	klog "github.com/go-kratos/kratos/v2/log"
	uow2 "github.com/go-saas/uow"
	kitconf "github.com/go-saas/kit/pkg/conf"
	"github.com/go-saas/kit/pkg/dal"
	"github.com/go-saas/kit/pkg/event"
	"github.com/go-saas/kit/pkg/event/trace"
	"github.com/go-saas/kit/user/private/biz"
)

func NewEventServer(c *kitconf.Data, conn dal.ConnName, logger klog.Logger, uowMgr uow2.Manager, tenantSeed biz.TenantSeedEventHandler) *event.ConsumerFactoryServer {
	e := c.Endpoints.GetEventMergedDefault(string(conn))
	srv := event.NewConsumerFactoryServer(e)
	srv.Use(event.ConsumerRecover(event.WithLogger(logger)), trace.Receive(), event.Logging(logger), event.ConsumerUow(uowMgr))
	srv.Append(tenantSeed)
	return srv
}
