package biz

import (
	"context"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/pkg/blob"
	"github.com/goxiaoy/go-saas-kit/pkg/dal"
	"github.com/goxiaoy/go-saas-kit/pkg/event/event"
	"github.com/goxiaoy/uow"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	NewTenantUserCase,
	NewRemoteEventHandler,
	NewTenantReadyEventHandler,
	NewConfigConnStrGenerator,
)

const ConnName dal.ConnName = "saas"

func LogoBlob(ctx context.Context, factory blob.Factory) blob.Blob {
	return factory.Get(ctx, "saas", false)
}

type SaasEventHandler event.Handler

//NewRemoteEventHandler handler for remote event
func NewRemoteEventHandler(l klog.Logger, uowMgr uow.Manager, tenantReady TenantReadyEventHandler) SaasEventHandler {
	return SaasEventHandler(event.RecoverHandler(event.UowHandler(uowMgr, event.ChainHandler(event.Handler(tenantReady))), event.WithLogger(l)))
}
