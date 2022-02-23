package service

import (
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authz"
)

func NewAuthorizationOption() *authz.Option {
	return authz.NewAuthorizationOption()
}

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewAuthorizationOption, NewMenuService)
