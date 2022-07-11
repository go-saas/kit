package server

import (
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-saas/kit/event"
	"github.com/go-saas/kit/event/trace"
	kitconf "github.com/go-saas/kit/pkg/conf"
	"github.com/go-saas/kit/pkg/dal"
	saasbiz "github.com/go-saas/kit/saas/private/biz"
	ubiz "github.com/go-saas/kit/user/private/biz"
	uow2 "github.com/go-saas/uow"
)

func NewEventServer(
	c *kitconf.Data,
	conn dal.ConnName,
	logger klog.Logger,
	uowMgr uow2.Manager,
	tenantReady saasbiz.TenantReadyEventHandler,
	tenantSeed ubiz.TenantSeedEventHandler,
	ur ubiz.UserRoleChangeEventHandler,
) *event.ConsumerFactoryServer {
	e := c.Endpoints.GetEventMergedDefault(string(conn))
	srv := event.NewConsumerFactoryServer(e)
	srv.Use(event.ConsumerRecover(event.WithLogger(logger)), trace.Receive(), event.Logging(logger), event.ConsumerUow(uowMgr))
	srv.Append(tenantReady)
	srv.Append(tenantSeed)
	srv.Append(ur)
	return srv
}
