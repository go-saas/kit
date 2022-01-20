package service

import (
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authorization"
)

func NewAuthorizationOption(userRole *UserRoleContributor) *authorization.Option {
	return authorization.NewAuthorizationOption(userRole)
}

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	NewUserRoleContributor,
	NewAuthorizationOption,
	NewUserService,
	NewAccountService,
	NewAuthService,
	NewRoleServiceService,
	NewPermissionService)
