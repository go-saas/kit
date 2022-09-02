package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	_ "github.com/go-saas/kit/pkg/blob/memory"
	_ "github.com/go-saas/kit/pkg/blob/os"
	_ "github.com/go-saas/kit/pkg/blob/s3"
	conf2 "github.com/go-saas/kit/pkg/conf"
	kitdi "github.com/go-saas/kit/pkg/di"
	"github.com/go-saas/kit/realtime/private/biz"
	"github.com/go-saas/saas/gorm"
	g "gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = kitdi.NewSet(
	NewData,
	NewMigrate,
	NewNotificationRepo,
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
func NewData(c *conf2.Data, dbProvider gorm.DbProvider, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		logger.Log(log.LevelInfo, log.DefaultMessageKey, "closing the data resources")
	}
	return &Data{
		DbProvider: dbProvider,
	}, cleanup, nil
}
