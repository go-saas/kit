package biz

import (
	"github.com/go-saas/kit/event"
	"github.com/go-saas/kit/pkg/dal"
	kitdi "github.com/go-saas/kit/pkg/di"
	"github.com/goava/di"
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

	kitdi.NewProvider(NewTenantSeedEventHandler, di.As(new(event.ConsumerHandler))),
	kitdi.NewProvider(NewUserRoleChangeEventHandler, di.As(new(event.ConsumerHandler))),

	//job
	NewUserMigrationTaskHandler)

const UserAvatarPath = "user/avatar"
