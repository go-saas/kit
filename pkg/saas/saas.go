package saas

import (
	"context"
	"github.com/goxiaoy/go-saas/common"
	"github.com/goxiaoy/go-saas/seed"
)

// RunWithTenantCache  get tenant config and cached into context
func RunWithTenantCache(ctx context.Context, store common.TenantStore, f func(ctx context.Context) error) error {
	tenantConfigProvider := common.NewDefaultTenantConfigProvider(common.NewDefaultTenantResolver(), store)
	_, ctx, err := tenantConfigProvider.Get(ctx)
	if err != nil {
		return err
	}
	return f(ctx)
}

type SeedFunc func(ctx context.Context, sCtx *seed.Context) error

func (s SeedFunc) Seed(ctx context.Context, sCtx *seed.Context) error {
	return s(ctx, sCtx)
}

func SeedChangeTenant(store common.TenantStore, next ...seed.Contributor) seed.Contributor {
	return SeedFunc(func(ctx context.Context, sCtx *seed.Context) error {
		return RunWithTenantCache(ctx, store, func(ctx context.Context) error {
			return seed.Chain(next...).Seed(ctx, sCtx)
		})
	})
}
