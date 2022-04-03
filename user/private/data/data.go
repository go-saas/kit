package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/goxiaoy/go-eventbus"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/casbin"
	"github.com/goxiaoy/go-saas-kit/pkg/blob"
	_ "github.com/goxiaoy/go-saas-kit/pkg/blob/memory"
	_ "github.com/goxiaoy/go-saas-kit/pkg/blob/os"
	_ "github.com/goxiaoy/go-saas-kit/pkg/blob/s3"
	"github.com/goxiaoy/go-saas-kit/pkg/email"
	event2 "github.com/goxiaoy/go-saas-kit/pkg/event"
	"github.com/goxiaoy/go-saas-kit/pkg/event/event"
	kitgorm "github.com/goxiaoy/go-saas-kit/pkg/gorm"
	"github.com/goxiaoy/go-saas-kit/pkg/lazy"
	kitredis "github.com/goxiaoy/go-saas-kit/pkg/redis"
	uow2 "github.com/goxiaoy/go-saas-kit/pkg/uow"
	"github.com/goxiaoy/go-saas-kit/user/private/biz"
	"github.com/goxiaoy/go-saas-kit/user/private/conf"
	"github.com/goxiaoy/go-saas/common"
	"github.com/goxiaoy/go-saas/data"
	"github.com/goxiaoy/go-saas/gorm"
	mail "github.com/xhit/go-simple-mail/v2"
	g "gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewConnStrResolver,
	kitgorm.NewDbProvider,
	kitgorm.NewDbOpener,
	wire.Value(eventbus.Default),
	NewRemoteEventReceiver,
	NewEventSender,
	uow2.NewUowManager,
	NewBlobFactory,
	NewRedis,
	kitredis.NewCache,
	NewEmailer,
	NewEnforcerProvider,
	NewUserRepo,
	NewRefreshTokenRepo,
	NewRoleRepo,
	NewUserTenantRepo,
	NewMigrate,
	NewUserTokenRepo,
	NewUserSettingRepo,
	NewUserAddrRepo)

const ConnName = "user"

// Data .
type Data struct {
	DbProvider gorm.DbProvider
}

func GetDb(ctx context.Context, provider gorm.DbProvider) *g.DB {
	db := provider.Get(ctx, ConnName)
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

func NewConnStrResolver(c *conf.Data, ts common.TenantStore) data.ConnStrResolver {
	return kitgorm.NewConnStrResolver(c.Endpoints, ts)
}

func NewEnforcerProvider(logger log.Logger, dbProvider gorm.DbProvider) (*casbin.EnforcerProvider, error) {
	return casbin.NewEnforcerProvider(logger, dbProvider, ConnName)
}

func NewBlobFactory(c *conf.Data) blob.Factory {
	return blob.NewFactory(c.Blobs)
}

func NewRedis(c *conf.Data) *redis.Client {
	if c.Endpoints.Redis == nil {
		panic("redis endpoints required")
	}
	r := kitredis.ResolveRedisEndpointOrDefault(c.Endpoints.Redis, ConnName)
	return kitredis.NewRedisClient(r)
}

func NewEmailer(c *conf.Data) *lazy.Of[*mail.SMTPClient] {
	return email.NewLazyClient(c.Endpoints)
}

func NewEventSender(c *conf.Data, logger log.Logger) (event.Sender, func(), error) {
	e := c.Endpoints.GetEventOrDefault(ConnName)
	return event2.NewEventSender(e, logger)
}

func NewRemoteEventReceiver(c *conf.Data, logger log.Logger, handler event.Handler) (event.Receiver, func(), error) {
	e := c.Endpoints.GetEventOrDefault(ConnName)
	return event2.NewEventReceiver(e, handler, logger)
}
