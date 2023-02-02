package gorm

import (
	"context"
	"gorm.io/gorm"
)

type (
	contextDbKey string
)

func NewContext(ctx context.Context, key string, db *gorm.DB) context.Context {
	return context.WithValue(ctx, contextDbKey(key), db)
}

func fromContext(ctx context.Context, key contextDbKey) (*gorm.DB, bool) {
	v, ok := ctx.Value(key).(*gorm.DB)
	return v, ok
}
