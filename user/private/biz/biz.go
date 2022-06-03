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
	NewUserManager,
	NewSignInManager,
	NewUserValidator,
	NewRoleManager,
	NewLookupNormalizer,
	NewRemoteEventHandler,
	NewEmailTokenProvider,
	NewPhoneTokenProvider,
	NewTwoStepTokenProvider,
	NewPasswordHasher,
	NewPasswordValidator,
	NewRoleSeed,
	NewUserSeed,
	NewPermissionSeeder,
	NewEmailSender,
	NewTenantSeedEventHandler)

func ProfileBlob(ctx context.Context, factory blob.Factory) blob.Blob {
	return factory.Get(ctx, "user", false)
}

//NewRemoteEventHandler handler for remote event
func NewRemoteEventHandler(l klog.Logger, uowMgr uow.Manager, tenantSeed TenantSeedEventHandler) event.Handler {
	return event.RecoverHandler(event.UowHandler(uowMgr, event.ChainHandler(event.Handler(tenantSeed))), event.WithLogger(l))
}
