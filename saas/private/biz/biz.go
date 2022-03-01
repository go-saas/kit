package biz

import (
	"context"
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/pkg/blob"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewTenantUserCase)

func LogoBlob(ctx context.Context, factory blob.Factory) blob.Blob {
	return factory.Get(ctx, "saas", false)
}
