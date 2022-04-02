package biz

import (
	"context"
	"github.com/google/wire"
	"github.com/goxiaoy/go-eventbus"
	"github.com/goxiaoy/go-saas-kit/pkg/blob"
	data2 "github.com/goxiaoy/go-saas-kit/pkg/data"
	"github.com/goxiaoy/go-saas-kit/pkg/event/event"
	v1 "github.com/goxiaoy/go-saas-kit/saas/event/v1"
	"github.com/goxiaoy/uow"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	NewTenantUserCase,
	NewLocalEventHook,
	NewRemoteEventHandler,
	NewTenantReadyEventHandler,
)

func LogoBlob(ctx context.Context, factory blob.Factory) blob.Blob {
	return factory.Get(ctx, "saas", false)
}

type EventHook interface {
}

//NewLocalEventHook hook with local event
func NewLocalEventHook(sender event.Sender) (EventHook, func(), error) {
	var cleanup = func() {
	}
	dispose1, err := eventbus.Subscribe[*data2.AfterCreate[*Tenant]]()(func(ctx context.Context, data *data2.AfterCreate[*Tenant]) error {
		event, err := event.NewMessageFromProto(&v1.TenantCreatedEvent{
			Id:         data.Entity.ID.String(),
			Name:       data.Entity.Name,
			Region:     data.Entity.Region,
			SeparateDb: data.Entity.SeparateDb,
		})
		if err != nil {
			return err
		}
		return sender.Send(ctx, event)
	})
	if err != nil {
		return nil, cleanup, err
	}

	return eventbus.Default, func() {
		dispose1.Dispose()
	}, nil
}

//NewRemoteEventHandler handler for remote event
func NewRemoteEventHandler(uowMgr uow.Manager, tenantReady TenantReadyEventHandler) event.Handler {
	return event.UowHandler(uowMgr, event.ChainHandler(event.Handler(tenantReady)))
}
