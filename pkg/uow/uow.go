package uow

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	gorm3 "github.com/goxiaoy/go-saas-kit/pkg/gorm"
	"github.com/goxiaoy/uow"
	gorm2 "github.com/goxiaoy/uow/gorm"
	"strings"
)

const (
	gormKind = "gorm"
)

func NewUowManager(cache *gorm3.DbCache) uow.Manager {
	return uow.NewManager(func(ctx context.Context, keys ...string) uow.TransactionalDb {
		kind := keys[0]
		key := keys[1]
		connStr := keys[2]
		if kind == gormKind {
			db, err := cache.GetOrSet(ctx, key, connStr)
			if err != nil {
				panic(err)
			}
			return gorm2.NewTransactionDb(db)
		}
		panic(errors.New(fmt.Sprintf("can not resolve %s", key)))
	})
}

var (
	safeMethods = []string{"GET", "HEAD", "OPTIONS", "TRACE"}
)

func contains(vals []string, s string) bool {
	for _, v := range vals {
		if v == s {
			return true
		}
	}

	return false
}

func Uow(l log.Logger, um uow.Manager) middleware.Middleware {
	logger := log.NewHelper(log.With(l, "module", "uow"))
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			var res interface{}
			var err error

			if t, ok := transport.FromServerContext(ctx); ok {
				//resolve by operation
				if len(t.Operation()) > 0 && skipOperation(t.Operation()) {
					//skip unit of work
					logger.Debugf("safe operation %s. skip uow", t.Operation())
					return handler(ctx, req)
				}
				// can not identify
				if ht, ok := t.(*khttp.Transport); ok {
					if contains(safeMethods, ht.Request().Method) {
						//safe method skip unit of work
						logger.Debugf("safe method %s. skip uow", ht.Request().Method)
						return handler(ctx, req)
					}
				}
			}

			// wrap into new unit of work
			logger.Debugf("run into unit of work")
			err = um.WithNew(ctx, func(ctx context.Context) error {
				var err error
				res, err = handler(ctx, req)
				return err
			})
			return res, err
		}
	}
}

//useOperation return true if operation action not start with "get" and "list" (case-insensitive)
func skipOperation(operation string) bool {
	s := strings.Split(operation, "/")
	act := strings.ToLower(s[len(s)-1])
	return strings.HasPrefix(act, "get") || strings.HasPrefix(act, "list")
}
