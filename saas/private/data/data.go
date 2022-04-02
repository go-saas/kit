package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/goxiaoy/go-eventbus"
	"github.com/goxiaoy/go-saas-kit/pkg/blob"
	event2 "github.com/goxiaoy/go-saas-kit/pkg/event"
	"github.com/goxiaoy/go-saas-kit/pkg/event/event"
	kitgorm "github.com/goxiaoy/go-saas-kit/pkg/gorm"
	uow2 "github.com/goxiaoy/go-saas-kit/pkg/uow"
	"github.com/goxiaoy/go-saas-kit/saas/private/conf"
	"github.com/goxiaoy/go-saas/common"
	"github.com/goxiaoy/go-saas/data"
	"github.com/goxiaoy/go-saas/gorm"
	g "gorm.io/gorm"

	_ "github.com/goxiaoy/go-saas-kit/pkg/blob/memory"
	_ "github.com/goxiaoy/go-saas-kit/pkg/blob/os"
	_ "github.com/goxiaoy/go-saas-kit/pkg/blob/s3"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewConnStrResolver,
	kitgorm.NewDbOpener,
	kitgorm.NewDbProvider,
	//local event bus
	wire.Value(eventbus.Default),
	NewRemoteEventReceiver,
	uow2.NewUowManager,
	NewTenantStore,
	NewBlobFactory,
	NewTenantRepo,
	NewEventSender,
	NewMigrate,
)

const ConnName = "saas"

// Data .
type Data struct {
	DbProvider gorm.DbProvider
}

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

func NewConnStrResolver(c *conf.Data, ts common.TenantStore) data.ConnStrResolver {
	return kitgorm.NewConnStrResolver(c.Endpoints, ts)
}
func NewBlobFactory(c *conf.Data) blob.Factory {
	return blob.NewFactory(c.Blobs)
}

func NewEventSender(c *conf.Data, logger log.Logger) (event.Sender, func(), error) {
	e := c.Endpoints.GetEventOrDefault(ConnName)
	return event2.NewEventSender(e, logger)
}

func NewRemoteEventReceiver(c *conf.Data, logger log.Logger, handler event.Handler) (event.Receiver, func(), error) {
	e := c.Endpoints.GetEventOrDefault(ConnName)
	return event2.NewEventReceiver(e, handler, logger)
}
