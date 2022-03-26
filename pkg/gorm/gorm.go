package gorm

import (
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	"github.com/goxiaoy/go-saas/common"
	"github.com/goxiaoy/go-saas/data"
	sgorm "github.com/goxiaoy/go-saas/gorm"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/gorm"
)

func NewDbOpener() (sgorm.DbOpener, func()) {
	return sgorm.NewDbOpener(func(db *gorm.DB) *gorm.DB {
		RegisterCallbacks(db)
		if err := db.Use(otelgorm.NewPlugin()); err != nil {
			panic(err)
		}
		return db
	})
}

func NewConnStrResolver(c *conf.Endpoints, ts common.TenantStore) data.ConnStrResolver {
	conn := make(data.ConnStrings, 1)
	for k, v := range c.Databases {
		conn[k] = v.Source
	}
	mr := common.NewMultiTenancyConnStrResolver(func() common.TenantStore {
		return common.NewCachedTenantStore(ts)
	}, data.NewConnStrOption(conn))
	return mr
}

func NewDbProvider(cs data.ConnStrResolver, c *sgorm.Config, opener sgorm.DbOpener) sgorm.DbProvider {
	return sgorm.NewDefaultDbProvider(cs, c, opener)
}
