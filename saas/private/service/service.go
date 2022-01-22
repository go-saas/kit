package service

import (
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authorization"
)

func NewAuthorizationOption() *authorization.Option {
	return authorization.NewAuthorizationOption()
}

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewTenantService, NewAuthorizationOption)
