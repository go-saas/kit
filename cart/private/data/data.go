package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/goxiaoy/go-eventbus"
	"github.com/goxiaoy/go-saas-kit/pkg/blob"
	_ "github.com/goxiaoy/go-saas-kit/pkg/blob/memory"
	_ "github.com/goxiaoy/go-saas-kit/pkg/blob/os"
	_ "github.com/goxiaoy/go-saas-kit/pkg/blob/s3"
	kitgorm "github.com/goxiaoy/go-saas-kit/pkg/gorm"
	suow "github.com/goxiaoy/go-saas-kit/pkg/uow"
	"github.com/goxiaoy/go-saas/common"
	"github.com/goxiaoy/go-saas/data"
	"github.com/goxiaoy/go-saas/gorm"
	"cart/private/conf"
	g "gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	kitgorm.NewDbOpener,
	wire.Value(eventbus.Default),
	kitgorm.NewDbProvider,
	NewConnStrResolver,
	suow.NewUowManager,
	NewBlobFactory,
	NewMigrate,
	NewPostRepo,
)

const ConnName = "cart"

// Data .
type Data struct {
	DbProvider gorm.DbProvider
}

func GetDb(ctx context.Context, provider gorm.DbProvider) *g.DB {
	db := provider.Get(ctx, ConnName)
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

func NewConnStrResolver(c *conf.Data, ts common.TenantStore) data.ConnStrResolver {
	return kitgorm.NewConnStrResolver(c.Endpoints, ts)
}

func NewBlobFactory(c *conf.Data) blob.Factory {
	return blob.NewFactory(c.Blobs)
}
