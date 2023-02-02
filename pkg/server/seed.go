package server

import (
	"context"
	uow2 "github.com/go-saas/kit/pkg/uow"
	"github.com/go-saas/saas"
	"github.com/go-saas/saas/seed"
	"github.com/go-saas/uow"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

// RunWithTenantCache  get tenant config and cached into context
func RunWithTenantCache(ctx context.Context, store saas.TenantStore, f func(ctx context.Context) error) error {
	tenantConfigProvider := saas.NewDefaultTenantConfigProvider(saas.NewDefaultTenantResolver(), store)
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

func SeedChangeTenant(store saas.TenantStore, next ...seed.Contrib) seed.Contrib {
	return SeedFunc(func(ctx context.Context, sCtx *seed.Context) error {
		return RunWithTenantCache(ctx, store, func(ctx context.Context) error {
			return seed.Chain(next...).Seed(ctx, sCtx)
		})
	})
}

func NewTraceContrib(next ...seed.Contrib) seed.Contrib {
	tracer := otel.Tracer("seeder")
	return SeedFunc(func(ctx context.Context, sCtx *seed.Context) (err error) {
		ctx, span := tracer.Start(ctx,
			"seed",
		)
		defer func() {
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			} else {
				span.SetStatus(codes.Ok, "OK")
			}
			span.End()
		}()
		err = seed.Chain(next...).Seed(ctx, sCtx)
		return
	})
}

func NewUowContrib(uow uow.Manager, next ...seed.Contrib) seed.Contrib {
	return SeedFunc(func(ctx context.Context, sCtx *seed.Context) error {
		if uow2.SkipFromContext(ctx, false) {
			return seed.Chain(next...).Seed(ctx, sCtx)
		}
		return uow.WithNew(ctx, func(ctx context.Context) error {
			return seed.Chain(next...).Seed(ctx, sCtx)
		})
	})
}
