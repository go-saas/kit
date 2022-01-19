package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	uow2 "github.com/goxiaoy/go-saas-kit/pkg/uow"
	"github.com/goxiaoy/go-saas-kit/saas/private/conf"
	"github.com/goxiaoy/go-saas/common"
	"github.com/goxiaoy/go-saas/data"
	"github.com/goxiaoy/go-saas/gorm"
	g "gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, gorm.NewDbOpener, uow2.NewUowManager, NewTenantStore, NewProvider, NewTenantRepo, NewMigrate)

const ConnName = "saas"

// Data .
type Data struct {
	DbProvider gorm.DbProvider
}

// GlobalData TODO better way to prevent cycle dependency
var GlobalData *Data

func GetDb(ctx context.Context, provider gorm.DbProvider) *g.DB {
	db := provider.Get(ctx, ConnName)
	return db
}

// NewData .
func NewData(c *conf.Data, dbProvider gorm.DbProvider, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		logger.Log(log.LevelInfo, "closing the data resources")
	}
	GlobalData = &Data{
		DbProvider: dbProvider,
	}
	return GlobalData, cleanup, nil
}

func NewProvider(c *conf.Data, cfg *gorm.Config, opener gorm.DbOpener, ts common.TenantStore, logger log.Logger) gorm.DbProvider {
	conn := make(data.ConnStrings, 1)
	for k, v := range c.Endpoints.Databases {
		conn[k] = v.Source
	}
	mr := common.NewMultiTenancyConnStrResolver(func() common.TenantStore {
		return ts
	}, data.NewConnStrOption(conn))
	r := gorm.NewDefaultDbProvider(mr, cfg, opener)
	return r
}
