package biz

import (
	"context"
	"github.com/go-saas/kit/event"
	"github.com/go-saas/kit/pkg/blob"
	"github.com/go-saas/kit/pkg/dal"
	kitdi "github.com/go-saas/kit/pkg/di"
	"github.com/goava/di"
)

// ProviderSet is biz providers.
var ProviderSet = kitdi.NewSet(
	NewTenantUserCase,
	kitdi.NewProvider(NewTenantReadyEventHandler, di.As(new(event.ConsumerHandler))),
	NewConfigConnStrGenerator,
)

const ConnName dal.ConnName = "saas"

func LogoBlob(ctx context.Context, factory blob.Factory) blob.Blob {
	return factory.Get(ctx, "saas", false)
}
