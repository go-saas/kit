package biz

import (
	"context"
	"github.com/go-saas/saas/seed"
)

type PostSeeder struct {
}

var _ seed.Contrib = (*PostSeeder)(nil)

func NewPostSeeder() *PostSeeder {
	return &PostSeeder{}
}

func (p *PostSeeder) Seed(ctx context.Context, sCtx *seed.Context) error {
	return nil
}
