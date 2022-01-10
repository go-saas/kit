package biz

import (
	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	NewUserManager,
	NewUserValidator,
	NewRoleManager,
	NewLookupNormalizer,
	NewPasswordHasher,
	NewPasswordValidator,
	NewRoleSeed,
	NewUserSeed,
	NewPermissionSeeder,
	NewAuthbossStoreWrapper)
