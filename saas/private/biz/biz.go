package biz

import (
	"context"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/pkg/blob"
	"github.com/goxiaoy/go-saas-kit/pkg/event/event"
	"github.com/goxiaoy/uow"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	NewTenantUserCase,
	NewLocalEventHook,
	NewRemoteEventHandler,
	NewTenantReadyEventHandler,
	NewConfigConnStrGenerator,
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
	return nil, cleanup, nil
}

type SaasEventHandler event.Handler

//NewRemoteEventHandler handler for remote event
func NewRemoteEventHandler(l klog.Logger, uowMgr uow.Manager, tenantReady TenantReadyEventHandler) SaasEventHandler {
	return SaasEventHandler(event.RecoverHandler(event.UowHandler(uowMgr, event.ChainHandler(event.Handler(tenantReady))), event.WithLogger(l)))
}
