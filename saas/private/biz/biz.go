package biz

import (
	"context"
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/pkg/blob"
	"github.com/goxiaoy/go-saas-kit/pkg/dal"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	NewTenantUserCase,
	NewTenantReadyEventHandler,
	NewConfigConnStrGenerator,
)

const ConnName dal.ConnName = "saas"

func LogoBlob(ctx context.Context, factory blob.Factory) blob.Blob {
	return factory.Get(ctx, "saas", false)
}
