package biz

import (
	"context"
	"github.com/go-saas/kit/pkg/blob"
	"github.com/go-saas/kit/pkg/dal"
	kitdi "github.com/go-saas/kit/pkg/di"
)

// ProviderSet is biz providers.
var ProviderSet = kitdi.NewSet(
	NewTenantUserCase,
	NewTenantReadyEventHandler,
	NewConfigConnStrGenerator,
)

const ConnName dal.ConnName = "saas"

func LogoBlob(ctx context.Context, factory blob.Factory) blob.Blob {
	return factory.Get(ctx, "saas", false)
}
