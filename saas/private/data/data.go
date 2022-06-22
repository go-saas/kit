package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	kconf "github.com/goxiaoy/go-saas-kit/pkg/conf"
	"github.com/goxiaoy/go-saas-kit/pkg/dal"
	"github.com/goxiaoy/go-saas-kit/saas/api"
	"github.com/goxiaoy/go-saas-kit/saas/private/biz"
	g "gorm.io/gorm"

	_ "github.com/goxiaoy/go-saas-kit/pkg/blob/memory"
	_ "github.com/goxiaoy/go-saas-kit/pkg/blob/os"
	_ "github.com/goxiaoy/go-saas-kit/pkg/blob/s3"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewTenantRepo,
	NewMigrate,
	api.NewTenantStore,
)

// Data .
type Data struct {
	DbProvider dal.ConstDbProvider
}

func GetDb(ctx context.Context, provider dal.ConstDbProvider) *g.DB {
	db := provider.Get(ctx, string(biz.ConnName))
	return db
}

// NewData .
func NewData(c *kconf.Data, dbProvider dal.ConstDbProvider, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		logger.Log(log.LevelInfo, "closing the data resources")
	}
	return &Data{
		DbProvider: dbProvider,
	}, cleanup, nil
}
