package job

import (
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/hibiken/asynq"
)

var DefaultProviderSet = wire.NewSet(NewAsynqClientOpt, NewAsynqClient)

type RedisFunc func() interface{}

func (r RedisFunc) MakeRedisClient() interface{} {
	return r()
}

func NewAsynqClientOpt(r *redis.Client) asynq.RedisConnOpt {
	return RedisFunc(func() interface{} {
		return r
	})
}

func NewAsynqClient(opt asynq.RedisConnOpt) (*asynq.Client, func()) {
	client := asynq.NewClient(opt)
	return client, func() {
		client.Close()
	}
}
