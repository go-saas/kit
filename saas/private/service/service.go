package service

import (
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/authorization/authorization"
	"github.com/goxiaoy/go-saas-kit/user/api"
)


func NewAuthorizationOption(userRole *api.RemoteRoleContributor) *authorization.Option {
	return authorization.NewAuthorizationOption(userRole)
}

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewTenantService,NewAuthorizationOption)
