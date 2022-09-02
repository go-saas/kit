package redis

import (
	"github.com/go-redis/redis/extra/redisotel/v8"
	"github.com/go-redis/redis/v8"
)

const defaultKey = "default"

func ResolveRedisEndpointOrDefault(endpoints map[string]*Config, key string) (*redis.UniversalOptions, error) {
	var opt *Config
	opt, ok := endpoints[key]
	if !ok {
		opt = endpoints[defaultKey]
	}
	redisOpt := &redis.UniversalOptions{
		Addrs: opt.Addrs,
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
	if opt.MasterName != nil {
		redisOpt.MasterName = *opt.MasterName
	}
	return redisOpt, nil
}

func NewRedisClient(r *redis.UniversalOptions) redis.UniversalClient {
	rdb := redis.NewUniversalClient(r)
	rdb.AddHook(redisotel.NewTracingHook())
	return rdb
}
