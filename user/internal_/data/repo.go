package data

import (
	"context"
	"github.com/goxiaoy/go-saas/gorm"
	g "gorm.io/gorm"
	gorm2 "github.com/goxiaoy/go-saas-kit/pkg/gorm"
)

const ConnKey = "User"

type Repo struct {
	gorm2.Repo
	DbProvider gorm.DbProvider
}

func (r *Repo) GetDb(ctx context.Context) *g.DB {
	return GetDb(ctx, r.DbProvider)
}
