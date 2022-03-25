package data

import (
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
)

const defaultKey = "default"

func ResolveRedisEndpointOrDefault(endpoints map[string]*conf.Redis, key string) *redis.Options {
	//TODO cluster
	var opt *conf.Redis
	opt, ok := endpoints[key]
	if !ok {
		opt = endpoints[defaultKey]
	}
	redisOpt := &redis.Options{
		Addr: opt.Addr,
	}

	if opt.Network != nil {
		redisOpt.Network = opt.Network.Value
	}
	if opt.Username != nil {
		redisOpt.Username = opt.Username.Value
	}
	if opt.Password != nil {
		redisOpt.Password = opt.Password.Value
	}
	if opt.Db != nil {
		redisOpt.DB = int(opt.Db.Value)
	}
	if opt.MaxRetries != nil {
		redisOpt.MaxRetries = int(opt.MaxRetries.Value)
	}
	if opt.MaxRetryBackoff != nil {
		redisOpt.MaxRetryBackoff = opt.MaxRetryBackoff.AsDuration()
	}
	if opt.MinRetryBackoff != nil {
		redisOpt.MinRetryBackoff = opt.MinRetryBackoff.AsDuration()
	}
	if opt.DialTimeout != nil {
		redisOpt.DialTimeout = opt.DialTimeout.AsDuration()
	}
	if opt.ReadTimeout != nil {
		redisOpt.ReadTimeout = opt.ReadTimeout.AsDuration()
	}
	if opt.WriteTimeout != nil {
		redisOpt.WriteTimeout = opt.WriteTimeout.AsDuration()
	}
	return redisOpt
}

func NewCache(client *redis.Client) *cache.Cache {
	c := cache.New(&cache.Options{
		Redis: client,
	})
	return c

}
