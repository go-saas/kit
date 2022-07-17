package job

import (
	"github.com/go-redis/redis/v8"
	kitdi "github.com/go-saas/kit/pkg/di"
	"github.com/hibiken/asynq"
)

var DefaultProviderSet = kitdi.NewSet(NewAsynqClientOpt, NewAsynqClient)

type RedisFunc func() interface{}

func (r RedisFunc) MakeRedisClient() interface{} {
	return r()
}

func NewAsynqClientOpt(r redis.UniversalClient) asynq.RedisConnOpt {
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

type Handler struct {
	Pattern string
	asynq.Handler
}

func NewHandler(pattern string, handler asynq.Handler) *Handler {
	return &Handler{Pattern: pattern, Handler: handler}
}
func NewHandlerFunc(pattern string, handler asynq.HandlerFunc) *Handler {
	return &Handler{Pattern: pattern, Handler: handler}
}

func RegisterHandlers(srv *Server, handlers ...*Handler) {
	for _, handler := range handlers {
		srv.Handle(handler.Pattern, handler.Handler)
	}
}
