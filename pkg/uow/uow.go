package uow

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-saas/kit/event"
	kitgorm "github.com/go-saas/kit/pkg/gorm"
	event2 "github.com/go-saas/kit/pkg/uow/event"
	"github.com/go-saas/uow"
	"github.com/go-saas/uow/gorm"
)

const (
	GormKind = "gorm"
)

type (
	skipUowKey struct{}
)

func NewUowManager(cache *kitgorm.DbCache) uow.Manager {
	return uow.NewManager(func(ctx context.Context, keys ...string) (uow.TransactionalDb, error) {
		kind := keys[0]
		if kind == GormKind {
			//see pkg/gorm/gorm.go#L188
			key := keys[1]
			connStr := keys[2]
			db, err := cache.GetOrSet(ctx, key, connStr)
			if err != nil {
				panic(err)
			}
			return gorm.NewTransactionDb(db), nil
		}
		if kind == event2.UowKind {
			if producer, ok := event.FromProducerContext(ctx); !ok || producer == nil {
				panic(errors.New("can not find producer"))
			} else {
				return event2.NewTransactional(ctx, producer), nil
			}
		}
		panic(errors.New(fmt.Sprintf("can not resolve %s", keys)))
	})
}

func SkipUow(ctx context.Context, skipOpt ...bool) context.Context {
	skip := true
	if len(skipOpt) > 0 {
		skip = skipOpt[0]
	}
	return context.WithValue(ctx, skipUowKey{}, skip)
}

func SkipFromContext(ctx context.Context, def bool) bool {
	if v, ok := ctx.Value(skipUowKey{}).(bool); ok {
		return v
	}
	return def
}
