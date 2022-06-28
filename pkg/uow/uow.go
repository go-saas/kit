package uow

import (
	"context"
	"errors"
	"fmt"
	gorm3 "github.com/go-saas/kit/pkg/gorm"
	"github.com/go-saas/uow"
	gorm2 "github.com/go-saas/uow/gorm"
)

const (
	gormKind = "gorm"
)

func NewUowManager(cache *gorm3.DbCache) uow.Manager {
	return uow.NewManager(func(ctx context.Context, keys ...string) (uow.TransactionalDb, error) {
		kind := keys[0]
		key := keys[1]
		connStr := keys[2]
		if kind == gormKind {
			db, err := cache.GetOrSet(ctx, key, connStr)
			if err != nil {
				panic(err)
			}
			return gorm2.NewTransactionDb(db), nil
		}
		panic(errors.New(fmt.Sprintf("can not resolve %s", keys)))
	})
}
