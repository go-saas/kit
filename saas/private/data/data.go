package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	kconf "github.com/go-saas/kit/pkg/conf"
	"github.com/go-saas/kit/pkg/dal"
	kitdi "github.com/go-saas/kit/pkg/di"
	"github.com/go-saas/kit/saas/api"
	"github.com/go-saas/kit/saas/private/biz"
	g "gorm.io/gorm"

	_ "github.com/go-saas/kit/pkg/blob/memory"
	_ "github.com/go-saas/kit/pkg/blob/os"
	_ "github.com/go-saas/kit/pkg/blob/s3"
)

// ProviderSet is data providers.
var ProviderSet = kitdi.NewSet(
	NewData,
	NewTenantRepo,
	NewMigrate,
	NewPlanRepo,
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
		logger.Log(log.LevelInfo, log.DefaultMessageKey, "closing the data resources")
	}
	return &Data{
		DbProvider: dbProvider,
	}, cleanup, nil
}
