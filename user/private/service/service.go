package service

import (
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authz"
)

func NewAuthorizationOption(userRole *UserRoleContributor) *authz.Option {
	return authz.NewAuthorizationOption(userRole)
}

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	NewUserTenantContributor,
	NewUserRoleContributor,
	NewAuthorizationOption,
	NewUserService,
	NewAccountService,
	NewAuthService,
	NewRoleServiceService,
	NewPermissionService)
