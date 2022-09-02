package server

import (
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-saas/kit/event"
	"github.com/go-saas/kit/event/trace"
	kitconf "github.com/go-saas/kit/pkg/conf"
	"github.com/go-saas/kit/pkg/dal"
	uow2 "github.com/go-saas/uow"
	"github.com/goava/di"
)

func NewEventServer(
	c *kitconf.Data,
	conn dal.ConnName,
	logger klog.Logger,
	uowMgr uow2.Manager,
	handlers []event.ConsumerHandler,
	container *di.Container,
) *event.ConsumerFactoryServer {
	e := c.Endpoints.GetEventMergedDefault(string(conn))
	srv := event.NewConsumerFactoryServer(e, container)
	srv.Use(event.ConsumerRecover(event.WithLogger(logger)), trace.Receive(), event.Logging(logger), event.ConsumerUow(uowMgr))
	for _, handler := range handlers {
		srv.Append(handler)
	}
	return srv
}
