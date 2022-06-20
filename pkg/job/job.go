package job

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/pkg/dal"
	"github.com/goxiaoy/go-saas-kit/pkg/lazy"
	"github.com/hibiken/asynq"
)

var DefaultProviderSet = wire.NewSet(NewAsynqClientOpt, NewAsynqClient, NewAsynqServer)

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

type LazyAsynqServer lazy.Of[*asynq.Server]

func NewAsynqServer(opt asynq.RedisConnOpt, name dal.ConnName) LazyAsynqServer {
	return lazy.New(func(ctx context.Context) (*asynq.Server, error) {
		//TODO read from config
		return asynq.NewServer(
			opt,
			asynq.Config{
				// Specify how many concurrent workers to use
				Concurrency: 10,
				// Optionally specify multiple queues with different priority.
				Queues: map[string]int{
					string(name): 1,
				},
				BaseContext: func() context.Context {
					return ctx
				},
				// See the godoc for other configuration options
			},
		), nil
	})

}
