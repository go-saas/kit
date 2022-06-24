package server

import (
	klog "github.com/go-kratos/kratos/v2/log"
	kitconf "github.com/goxiaoy/go-saas-kit/pkg/conf"
	"github.com/goxiaoy/go-saas-kit/pkg/dal"
	"github.com/goxiaoy/go-saas-kit/pkg/event"
	"github.com/goxiaoy/go-saas-kit/pkg/event/trace"
	"github.com/goxiaoy/go-saas-kit/user/private/biz"
	uow2 "github.com/goxiaoy/uow"
)

func NewEventServer(c *kitconf.Data, conn dal.ConnName, logger klog.Logger, uowMgr uow2.Manager, tenantSeed biz.TenantSeedEventHandler) *event.FactoryServer {
	e := c.Endpoints.GetEventMergedDefault(string(conn))
	srv := event.NewFactoryServer(e)
	srv.Use(event.Recover(event.WithLogger(logger)), trace.Receive(), event.Logging(logger), event.Uow(uowMgr))
	srv.Append(tenantSeed)
	return srv
}
