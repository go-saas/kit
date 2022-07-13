package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	_ "github.com/go-saas/kit/pkg/blob/memory"
	_ "github.com/go-saas/kit/pkg/blob/os"
	_ "github.com/go-saas/kit/pkg/blob/s3"
	kconf "github.com/go-saas/kit/pkg/conf"
	"github.com/go-saas/kit/sys/private/biz"
	"github.com/go-saas/saas/gorm"
	"github.com/google/wire"
	g "gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewMigrate,
	NewMenuRepo,
)

// Data .
type Data struct {
	DbProvider gorm.DbProvider
}

func GetDb(ctx context.Context, provider gorm.DbProvider) *g.DB {
	db := provider.Get(ctx, string(biz.ConnName))
	return db
}

// NewData .
func NewData(c *kconf.Data, dbProvider gorm.DbProvider, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		logger.Log(log.LevelInfo, log.DefaultMessageKey, "closing the data resources")
	}
	return &Data{
		DbProvider: dbProvider,
	}, cleanup, nil
}
