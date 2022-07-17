package biz

import (
	"context"
	"github.com/go-saas/kit/pkg/blob"
	"github.com/go-saas/kit/pkg/dal"
	kitdi "github.com/go-saas/kit/pkg/di"
)

const ConnName dal.ConnName = "user"

// ProviderSet is biz providers.
var ProviderSet = kitdi.NewSet(
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
	NewUserRoleChangeEventHandler,
	//job
	NewUserMigrationTaskHandler)

func ProfileBlob(ctx context.Context, factory blob.Factory) blob.Blob {
	return factory.Get(ctx, "user", false)
}
