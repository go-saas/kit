package data

import (
	"context"
	"github.com/go-saas/kit/pkg/dal"
	"gorm.io/gorm"
)

func GetDb(ctx context.Context, provider dal.ConstDbProvider, connName dal.ConnName) *gorm.DB {
	return provider.Get(ctx, string(connName))
}
