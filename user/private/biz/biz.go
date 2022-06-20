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

const ConnName dal.ConnName = "user"

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	NewUserManager,
	NewSignInManager,
	NewUserValidator,
	NewRoleManager,
	NewLookupNormalizer,

	NewEmailTokenProvider,
	NewPhoneTokenProvider,
	NewTwoStepTokenProvider,
	NewPasswordHasher,
	NewPasswordValidator,
	NewRoleSeed,
	NewUserSeed,
	NewPermissionSeeder,
	NewEmailSender,
	//event
	NewRemoteEventHandler,
	NewTenantSeedEventHandler,
	//job
	NewUserMigrationTaskHandler)

func ProfileBlob(ctx context.Context, factory blob.Factory) blob.Blob {
	return factory.Get(ctx, "user", false)
}

type UserEventHandler event.Handler

//NewRemoteEventHandler handler for remote event
func NewRemoteEventHandler(l klog.Logger, uowMgr uow.Manager, tenantSeed TenantSeedEventHandler) UserEventHandler {
	return UserEventHandler(event.RecoverHandler(event.UowHandler(uowMgr, event.ChainHandler(event.Handler(tenantSeed))), event.WithLogger(l)))
}
