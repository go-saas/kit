package gorm

import (
	"context"
	"github.com/goxiaoy/go-saas-kit/pkg/uow"
	"gorm.io/gorm"
)

func NewCurrent(ctx context.Context, k uow.UnitOfWorkKey, db *gorm.DB) context.Context {
	return context.WithValue(ctx, k, db)
}

func FromContext(ctx context.Context, k uow.UnitOfWorkKey) (db *gorm.DB, ok bool) {
	db, ok = ctx.Value(k).(*gorm.DB)
	return
}

func ChangeCurrent(ctx context.Context, k uow.UnitOfWorkKey, db *gorm.DB) (context.Context, uow.CancelFunc) {
	current, _ := FromContext(ctx, k)
	newCtx := NewCurrent(ctx, k, db)
	return newCtx, func(ctx context.Context) context.Context {
		return NewCurrent(ctx, k, current)
	}
}
