package service

import (
	v1 "github.com/go-saas/kit/oidc/api/client/v1"
	v12 "github.com/go-saas/kit/oidc/api/key/v1"
	kitdi "github.com/go-saas/kit/pkg/di"
	"github.com/goava/di"
)

var ProviderSet = kitdi.NewSet(
	kitdi.NewProvider(NewClientService, di.As(new(v1.ClientServiceServer))),
	kitdi.NewProvider(NewKeyService, di.As(new(v12.KeyServiceServer))),
)
