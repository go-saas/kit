package uow

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	"github.com/goxiaoy/go-saas/data"
	"github.com/goxiaoy/go-saas/gorm"
	"github.com/goxiaoy/uow"
	gorm2 "github.com/goxiaoy/uow/gorm"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	g "gorm.io/gorm"
	"gorm.io/gorm/schema"
	"strings"
)

func NewGormConfig(databases *conf.Endpoints, name string) *gorm.Config {
	var c *conf.Database
	var ok bool
	c, ok = databases.Databases[name]
	if !ok {
		c, ok = databases.Databases[data.Default]
	}
	tp := ""
	if c.TablePrefix == nil {
		tp = fmt.Sprintf("kit_%s_", name)
	} else {
		tp = c.TablePrefix.Value
	}
	cfg := &gorm.Config{
		Debug: c.Debug,
		Cfg: &g.Config{NamingStrategy: schema.NamingStrategy{
			TablePrefix: tp,
		}},
	}
	if c.Driver == "mysql" {
		cfg.Dialect = func(s string) g.Dialector {
			return mysql.Open(s)
		}
	}
	if c.Driver == "sqlite" {
		cfg.Dialect = func(s string) g.Dialector {
			return sqlite.Open(s)
		}
		//https://github.com/go-gorm/gorm/issues/2875
		cfg.MaxOpenConn = 1
		cfg.MaxIdleConn = 1
	}
	return cfg
}

func NewUowManager(cfg *gorm.Config, config *uow.Config, opener gorm.DbOpener) uow.Manager {
	return uow.NewManager(config, func(ctx context.Context, kind, key string) uow.TransactionalDb {
		if kind == gorm.DbKind {
			db, err := opener.Open(cfg, key)
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
				if ht, ok := t.(*khttp.Transport); ok {
					if contains(safeMethods, ht.Request().Method) {
						//safe method skip unit of work
						logger.Debugf("safe method %s. skip uow", ht.Request().Method)
						return handler(ctx, req)
					}
				} else {
					//resolve by operation
					if !useOperation(t.Operation()) {
						//skip unit of work
						logger.Debugf("safe operation %s. skip uow", t.Operation())
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
func useOperation(operation string) bool {
	s := strings.Split(operation, "/")
	act := strings.ToLower(s[len(s)-1])
	return !strings.HasPrefix(act, "get") && !strings.HasPrefix(act, "list")
}
