package biz

import (
	"context"
	"github.com/go-saas/kit/pkg/blob"
	"github.com/go-saas/kit/pkg/dal"
	"github.com/google/wire"
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
	NewPasswordHasher,
	NewPasswordValidator,
	NewRoleSeed,
	NewUserSeed,
	NewPermissionSeeder,
	NewEmailSender,

	NewTenantSeedEventHandler,
	//job
	NewUserMigrationTaskHandler)

func ProfileBlob(ctx context.Context, factory blob.Factory) blob.Blob {
	return factory.Get(ctx, "user", false)
}
