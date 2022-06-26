package redis

import (
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/extra/redisotel/v8"
	"github.com/go-redis/redis/v8"
	"github.com/go-saas/kit/pkg/conf"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

const defaultKey = "default"

func ResolveRedisEndpointOrDefault(endpoints map[string]*conf.Redis, key string) (*redis.Options, error) {
	//TODO cluster
	var opt *conf.Redis
	opt, ok := endpoints[key]
	if !ok {
		opt = endpoints[defaultKey]
	}
	return redis.ParseURL(opt.Url)
}

func NewRedisClient(r *redis.Options) *redis.Client {
	rdb := redis.NewClient(r)
	rdb.AddHook(redisotel.NewTracingHook(redisotel.WithAttributes(semconv.NetPeerNameKey.String(r.Addr), semconv.NetPeerPortKey.String(r.Addr))))
	return rdb
}

func NewCache(client *redis.Client) *cache.Cache {
	c := cache.New(&cache.Options{
		Redis: client,
	})
	return c

}
