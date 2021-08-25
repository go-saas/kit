package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	uow2 "github.com/goxiaoy/go-saas-kit/pkg/uow"
	"github.com/goxiaoy/go-saas-kit/user/internal_/biz"
	"github.com/goxiaoy/go-saas-kit/user/internal_/conf"
	"github.com/goxiaoy/go-saas/common"
	"github.com/goxiaoy/go-saas/data"
	"github.com/goxiaoy/go-saas/gorm"
	"github.com/goxiaoy/uow"
	g "gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, gorm.NewDbOpener, uow2.NewUowManager, NewTenantStore, NewProvider, NewUserRepo, NewRefreshTokenRepo, NewRoleRepo, NewMigrate)

const ConnName = "user"

// Data .
type Data struct {
	DbProvider gorm.DbProvider
}

func GetDb(ctx context.Context, provider gorm.DbProvider) *g.DB {
	db := provider.Get(ctx, ConnKey)
	if err := db.SetupJoinTable(&biz.User{}, "Roles", &biz.UserRole{}); err != nil {
		panic(err)
	}
	return db
}

// NewData .
func NewData(c *conf.Data, dbProvider gorm.DbProvider, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		logger.Log(log.LevelInfo, "closing the data resources")
	}
	return &Data{
		DbProvider: dbProvider,
	}, cleanup, nil
}

// NewTenantStore TODO replace with correct tenant store
func NewTenantStore() common.TenantStore {
	return common.NewMemoryTenantStore(
		[]common.TenantConfig{})
}

func NewProvider(c *conf.Data, cfg *gorm.Config, opener gorm.DbOpener, uow uow.Manager, ts common.TenantStore, logger log.Logger) gorm.DbProvider {
	ct := common.ContextCurrentTenant{}

	conn := make(data.ConnStrings, 1)
	for k,v :=range c.Databases.Databases{
		conn[k]=v.Source
	}
	mr := common.NewMultiTenancyConnStrResolver(ct, func() common.TenantStore {
		return ts
	}, data.NewConnStrOption(conn))
	r := gorm.NewDefaultDbProvider(mr, cfg, uow, opener)
	return r
}
