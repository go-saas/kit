package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	redisotel "github.com/go-redis/redis/extra/redisotel/v8"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/casbin"
	"github.com/goxiaoy/go-saas-kit/pkg/blob"
	_ "github.com/goxiaoy/go-saas-kit/pkg/blob/memory"
	_ "github.com/goxiaoy/go-saas-kit/pkg/blob/os"
	_ "github.com/goxiaoy/go-saas-kit/pkg/blob/s3"
	data2 "github.com/goxiaoy/go-saas-kit/pkg/data"
	"github.com/goxiaoy/go-saas-kit/pkg/email"
	kitgorm "github.com/goxiaoy/go-saas-kit/pkg/gorm"
	"github.com/goxiaoy/go-saas-kit/pkg/lazy"
	uow2 "github.com/goxiaoy/go-saas-kit/pkg/uow"
	"github.com/goxiaoy/go-saas-kit/user/private/biz"
	"github.com/goxiaoy/go-saas-kit/user/private/conf"
	"github.com/goxiaoy/go-saas/common"
	"github.com/goxiaoy/go-saas/data"
	"github.com/goxiaoy/go-saas/gorm"
	mail "github.com/xhit/go-simple-mail/v2"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	g "gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	kitgorm.NewDbOpener,
	uow2.NewUowManager,
	NewBlobFactory,
	NewProvider,
	NewRedis,
	data2.NewCache,
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

func NewEnforcerProvider(dbProvider gorm.DbProvider) *casbin.EnforcerProvider {
	return casbin.NewEnforcerProvider(dbProvider, ConnName)
}

func NewBlobFactory(c *conf.Data) blob.Factory {
	return blob.NewFactory(c.Blobs)
}

func NewRedis(c *conf.Data) *redis.Client {
	if c.Endpoints.Redis == nil {
		panic("redis endpoints required")
	}
	r := data2.ResolveRedisEndpointOrDefault(c.Endpoints.Redis, ConnName)
	rdb := redis.NewClient(r)
	rdb.AddHook(redisotel.NewTracingHook(redisotel.WithAttributes(semconv.NetPeerNameKey.String(r.Addr), semconv.NetPeerPortKey.String(r.Addr))))
	return rdb
}

func NewEmailer(c *conf.Data) *lazy.Of[*mail.SMTPClient] {
	return email.NewLazyClient(c.Endpoints)
}
