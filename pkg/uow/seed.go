package uow

import (
	"context"
	"github.com/goxiaoy/go-saas/seed"
	"github.com/goxiaoy/uow"
)

type Contributor struct {
	uow uow.Manager
	up  seed.Contributor
}

var _ seed.Contributor = (*Contributor)(nil)

func NewUowContributor(uow uow.Manager, up seed.Contributor) *Contributor {
	return &Contributor{
		uow: uow,
		up:  up,
	}
}

func (u *Contributor) Seed(ctx context.Context, sCtx *seed.Context) error {
	return u.uow.WithNew(ctx, func(ctx context.Context) error {
		return u.up.Seed(ctx, sCtx)
	})
}
