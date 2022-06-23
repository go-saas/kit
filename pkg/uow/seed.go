package uow

import (
	"context"
	"github.com/goxiaoy/go-saas-kit/pkg/saas"
	"github.com/goxiaoy/go-saas/seed"
	"github.com/goxiaoy/uow"
)

func NewUowContrib(uow uow.Manager, next ...seed.Contrib) seed.Contrib {
	return saas.SeedFunc(func(ctx context.Context, sCtx *seed.Context) error {
		return uow.WithNew(ctx, func(ctx context.Context) error {
			return seed.Chain(next...).Seed(ctx, sCtx)
		})
	})
}
