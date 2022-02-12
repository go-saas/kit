package biz

import (
	"context"
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/pkg/blob"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	NewUserManager,
	NewSignInManager,
	NewUserValidator,
	NewRoleManager,
	NewLookupNormalizer,
	NewPasswordHasher,
	NewPasswordValidator,
	NewRoleSeed,
	NewUserSeed,
	NewPermissionSeeder)

func ProfileBlob(ctx context.Context, factory blob.Factory) blob.Blob {
	return factory.Get(ctx, "user", false)
}
