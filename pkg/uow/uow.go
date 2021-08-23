package uow

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/goxiaoy/go-saas/gorm"
	"github.com/goxiaoy/uow"
	gorm2 "github.com/goxiaoy/uow/gorm"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	g "gorm.io/gorm"
)

func NewGormConfig(debug bool,driver string) *gorm.Config {
	cfg := &gorm.Config{
		Debug: debug,
		Cfg: &g.Config{},
	}
	if driver == "mysql"{
		cfg.Dialect = func(s string) g.Dialector {
			return mysql.Open(s)
		}
	}
	if driver == "sqlite"{
		cfg.Dialect = func(s string) g.Dialector {
			return sqlite.Open(s)
		}
		//https://github.com/go-gorm/gorm/issues/2875
		cfg.MaxOpenConn=1
		cfg.MaxIdleConn=1
	}
	return cfg
}

func NewUowManager(cfg *gorm.Config,config *uow.Config, opener gorm.DbOpener) uow.Manager {
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

func Uow(l log.Logger,um uow.Manager) middleware.Middleware {
	logger := log.NewHelper(l)
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			var res interface{}
			var err error
			// wrap into new unit of work
			logger.Debugf("run into unit of work")
			err = um.WithNew(ctx,func (ctx context.Context) error{
				var err error
				 res,err = handler(ctx, req)
				 return err
			})
			return res,err
		}
	}
}