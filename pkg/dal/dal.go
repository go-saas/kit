package dal

import (
	"github.com/eko/gocache/v3/cache"
	"github.com/eko/gocache/v3/store"
	"github.com/go-redis/redis/v8"
	"github.com/go-saas/kit/event"
	"github.com/go-saas/kit/event/trace"
	"github.com/go-saas/kit/pkg/blob"
	kitconf "github.com/go-saas/kit/pkg/conf"
	kitdi "github.com/go-saas/kit/pkg/di"
	"github.com/go-saas/kit/pkg/email"
	kitgorm "github.com/go-saas/kit/pkg/gorm"
	kitredis "github.com/go-saas/kit/pkg/redis"
	kituow "github.com/go-saas/kit/pkg/uow"
	"github.com/go-saas/saas"
	"github.com/go-saas/saas/data"
	sgorm "github.com/go-saas/saas/gorm"
	"github.com/goava/di"
	"github.com/goxiaoy/go-eventbus"
)

type (
	ConnName        string
	ConstDbProvider sgorm.DbProvider
)

var (
	//DefaultProviderSet shared provider for all data layer
	DefaultProviderSet = kitdi.NewSet(
		NewConnStrResolver,
		NewConstantConnStrResolver,
		kitgorm.NewDbCache,

		kitgorm.NewDbProvider,
		NewConstDbProvider,

		kituow.NewUowManager,

		NewBlobFactory,

		NewRedis,
		NewCacheStore,
		kitdi.NewProvider(NewStringCacheManager, di.As(new(cache.CacheInterface[string]))),

		NewEmailer,
		NewEventSender,
		func() *eventbus.EventBus { return eventbus.Default },
	)
)

func NewConnStrResolver(c *kitconf.Data, ts saas.TenantStore) data.ConnStrResolver {
	return kitgorm.NewConnStrResolver(c.Endpoints, ts)
}

// NewConstantConnStrResolver ignore multi-tenancy
func NewConstantConnStrResolver(c *kitconf.Data) data.ConnStrings {
	// ignore multi-tenancy
	conn := make(data.ConnStrings, 1)
	for k, v := range c.Endpoints.Databases {
		conn[k] = v.Source
	}
	return conn
}

// NewConstDbProvider ignore multi-tenancy
func NewConstDbProvider(cache *kitgorm.DbCache, cs data.ConnStrings, d *kitconf.Data) ConstDbProvider {
	return kitgorm.NewDbProvider(cache, cs, d)
}

func NewBlobFactory(c *kitconf.Data) blob.Factory {
	return blob.NewFactory(c.Blobs)
}

func NewEmailer(c *kitconf.Data) email.LazyClient {
	return email.NewLazyClient(c.Endpoints)
}

func NewRedis(c *kitconf.Data, name ConnName) (redis.UniversalClient, error) {
	if c.Endpoints.Redis == nil {
		panic("redis endpoints required")
	}
	r, err := kitredis.ResolveRedisEndpointOrDefault(c.Endpoints.Redis, string(name))
	return kitredis.NewRedisClient(r), err
}

func NewCacheStore(client redis.UniversalClient) store.StoreInterface {
	return store.NewRedis(client)
}

func NewStringCacheManager(store store.StoreInterface) *cache.Cache[string] {
	return cache.New[string](store)
}

func NewEventSender(c *kitconf.Data, name ConnName) (event.Producer, func(), error) {
	e := c.Endpoints.GetEventMergedDefault(string(name))
	ret, err := event.NewFactoryProducer(e)
	if err != nil {
		return nil, func() {}, err
	}
	ret.Use(trace.Send())
	return ret, func() {
		ret.Close()
	}, nil

}
