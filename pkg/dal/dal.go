package dal

import (
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/goxiaoy/go-eventbus"
	"github.com/goxiaoy/go-saas"
	"github.com/goxiaoy/go-saas-kit/pkg/blob"
	kitconf "github.com/goxiaoy/go-saas-kit/pkg/conf"
	"github.com/goxiaoy/go-saas-kit/pkg/email"
	event "github.com/goxiaoy/go-saas-kit/pkg/event"
	"github.com/goxiaoy/go-saas-kit/pkg/event/trace"
	kitgorm "github.com/goxiaoy/go-saas-kit/pkg/gorm"
	kitredis "github.com/goxiaoy/go-saas-kit/pkg/redis"
	kituow "github.com/goxiaoy/go-saas-kit/pkg/uow"

	"github.com/goxiaoy/go-saas/data"
	sgorm "github.com/goxiaoy/go-saas/gorm"
)

type (
	ConnName        string
	ConstDbProvider sgorm.DbProvider
)

var (
	//DefaultProviderSet shared provider for all data layer
	DefaultProviderSet = wire.NewSet(
		NewConnStrResolver,
		NewConstantConnStrResolver,
		kitgorm.NewDbCache,

		kitgorm.NewDbProvider,
		NewConstDbProvider,

		kituow.NewUowManager,

		NewBlobFactory,

		NewRedis,
		kitredis.NewCache,

		NewEmailer,
		NewEventSender,
		wire.Value(eventbus.Default),
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

func NewRedis(c *kitconf.Data, name ConnName) (*redis.Client, error) {
	if c.Endpoints.Redis == nil {
		panic("redis endpoints required")
	}
	r, err := kitredis.ResolveRedisEndpointOrDefault(c.Endpoints.Redis, string(name))
	return kitredis.NewRedisClient(r), err
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
