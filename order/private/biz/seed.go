package biz

import (
	"context"
	"github.com/goxiaoy/go-saas/seed"
)

type PostSeeder struct {
}

var _ seed.Contributor = (*PostSeeder)(nil)

func NewPostSeeder() *PostSeeder {
	return &PostSeeder{}
}

func (p *PostSeeder) Seed(ctx context.Context, sCtx *seed.Context) error {
	return nil
}
