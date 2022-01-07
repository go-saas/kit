package data

import (
	"context"
	gorm2 "github.com/goxiaoy/go-saas-kit/pkg/gorm"
	"github.com/goxiaoy/go-saas/gorm"
	g "gorm.io/gorm"
)

type Repo struct {
	gorm2.Repo
	DbProvider gorm.DbProvider
}

func (r *Repo) GetDb(ctx context.Context) *g.DB {
	return GetDb(ctx, r.DbProvider)
}
