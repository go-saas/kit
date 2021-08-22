package data

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas/common"
	"github.com/goxiaoy/go-saas/data"
	"github.com/goxiaoy/go-saas/gorm"
	"github.com/goxiaoy/uow"
	gorm2 "github.com/goxiaoy/uow/gorm"
	"gorm.io/driver/mysql"
	g "gorm.io/gorm"
	"github.com/goxiaoy/go-saas-kit/user/internal_/biz"
	"github.com/goxiaoy/go-saas-kit/user/internal_/conf"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData,gorm.NewDbOpener,NewProvider,NewUserRepo,NewRefreshTokenRepo,NewRoleRepo,NewMigrate)

// Data .
type Data struct {
	DbProvider gorm.DbProvider
}

func GetDb(ctx context.Context, provider gorm.DbProvider) *g.DB {
	db :=provider.Get(ctx, ConnKey)
	if err :=db.SetupJoinTable(&biz.User{}, "Roles", &biz.UserRole{});err!=nil{
		panic(err)
	}
	return db
}

// NewData .
func NewData(c *conf.Data,dbProvider gorm.DbProvider, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		logger.Log(log.LevelInfo, "closing the data resources")
	}
	return &Data{
		DbProvider: dbProvider,
	}, cleanup, nil
}

func NewProvider(c *conf.Data, opener gorm.DbOpener, logger log.Logger) gorm.DbProvider{

	cfg := &gorm.Config{
		//Debug: true,
		Dialect: func(s string) g.Dialector {
			return mysql.Open(s)
		},
		Cfg: &g.Config{},
	}

	ct := common.ContextCurrentTenant{}
	//TODO replace with correct tenant store
	ts := common.NewMemoryTenantStore(
		[]common.TenantConfig{})
	conn := make(data.ConnStrings, 1)
	conn.SetDefault(c.Database.Source)
	mr := common.NewMultiTenancyConnStrResolver(ct, func() common.TenantStore {
		return ts
	}, data.NewConnStrOption(conn))
	uow := uow.NewManager(&uow.Config{SupportNestedTransaction: false}, func(ctx context.Context, kind, key string) uow.TransactionalDb {
		if kind == gorm.GormDbKind {
			db, err := opener.Open(cfg, key)
			if err != nil {
				panic(err)
			}
			return gorm2.NewTransactionDb(db)
		}
		panic(errors.New(fmt.Sprintf("can not resolve %s", key)))
	})
	r := gorm.NewDefaultDbProvider(mr, cfg, uow, opener)
	return r
}
