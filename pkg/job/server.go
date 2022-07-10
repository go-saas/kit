package job

import (
	"context"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-saas/lazy"
	"github.com/hibiken/asynq"
)

type Server struct {
	*asynq.ServeMux
	server LazyAsynqServer
}

var _ transport.Server = (*Server)(nil)

type ServerOption func(opt *asynq.Config)

func WithQueues(q map[string]int) ServerOption {
	return func(opt *asynq.Config) {
		opt.Queues = q
	}
}

func WithConcurrency(c int) ServerOption {
	return func(opt *asynq.Config) {
		opt.Concurrency = c
	}
}

func NewServer(opt asynq.RedisConnOpt, opts ...ServerOption) *Server {
	ret := &Server{server: newAsynqServer(opt, opts...), ServeMux: asynq.NewServeMux()}
	return ret
}

func Logging(logger klog.Logger) asynq.MiddlewareFunc {
	return func(handler asynq.Handler) asynq.Handler {
		return asynq.HandlerFunc(func(ctx context.Context, task *asynq.Task) error {
			err := handler.ProcessTask(ctx, task)
			if err != nil {
				_ = klog.WithContext(ctx, logger).Log(klog.LevelError,
					klog.DefaultMessageKey, err.Error(),
					"task", task.Type())
			} else {
				_ = klog.WithContext(ctx, logger).Log(klog.LevelInfo,
					"task", task.Type())
			}
			return err
		})
	}
}

func (s *Server) Start(ctx context.Context) error {
	asynqServer, err := s.server.Value(ctx)
	if err != nil {
		return err
	}
	return asynqServer.Start(s)
}

func (s *Server) Stop(ctx context.Context) error {
	asynqServer, err := s.server.Value(ctx)
	if err != nil {
		return err
	}
	asynqServer.Shutdown()
	return nil
}

type LazyAsynqServer lazy.Of[*asynq.Server]

func newAsynqServer(opt asynq.RedisConnOpt, opts ...ServerOption) LazyAsynqServer {
	return lazy.New(func(ctx context.Context) (*asynq.Server, error) {
		cfg := asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 10,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{},
			BaseContext: func() context.Context {
				return ctx
			},
			// See the godoc for other configuration options
		}
		for _, option := range opts {
			option(&cfg)
		}
		return asynq.NewServer(
			opt,
			cfg,
		), nil
	})

}
