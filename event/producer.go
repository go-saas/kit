package event

import (
	"context"
	"github.com/goava/di"
	"io"
	"sync"
)

type Producer interface {
	io.Closer
	Send(ctx context.Context, msg Event) error
	BatchSend(ctx context.Context, msg []Event) error
}

type ProducerMux struct {
	Producer
	mu  sync.RWMutex
	mws []ProducerMiddlewareFunc
}

func NewFactoryProducer(cfg *Config, container *di.Container) (*ProducerMux, error) {
	_typeProducerMux.RLock()
	defer _typeProducerMux.RUnlock()
	t, err := resolveType(cfg)
	if err != nil {
		panic(err)
	}
	if r, ok := _typeProducerRegister[t]; !ok {
		panic(cfg.Type + " event producer not registered")
	} else {
		return r(cfg, container)
	}
}

func (s *ProducerMux) Close() error {
	return s.Producer.Close()
}

type ProducerMiddlewareFunc func(HandlerOf[any]) HandlerOf[any]

func ChainProducer(m ...ProducerMiddlewareFunc) ProducerMiddlewareFunc {
	return func(next HandlerOf[any]) HandlerOf[any] {
		for i := len(m) - 1; i >= 0; i-- {
			next = m[i](next)
		}
		return next
	}
}

// NewProducer create a *SendWrap with middleware ability
func NewProducer(next Producer) *ProducerMux {
	ret := &ProducerMux{
		Producer: next,
	}
	ret.Use(ret.wrapContext())
	return ret
}

func (s *ProducerMux) Use(m ...ProducerMiddlewareFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, mf := range m {
		s.mws = append(s.mws, mf)
	}
}

func (s *ProducerMux) Send(ctx context.Context, msg Event) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c := ChainProducer(s.mws...)
	h := c(HandlerFuncOf[any](func(ctx context.Context, e any) error {
		// put the real send as inner
		return s.Producer.Send(ctx, e.(Event))
	}))
	return h.Process(ctx, msg)
}

func (s *ProducerMux) BatchSend(ctx context.Context, msg []Event) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	c := ChainProducer(s.mws...)
	h := c(HandlerFuncOf[any](func(ctx context.Context, e any) error {
		return s.Producer.BatchSend(ctx, e.([]Event))
	}))
	return h.Process(ctx, msg)
}

func (s *ProducerMux) wrapContext() ProducerMiddlewareFunc {
	return func(next HandlerOf[any]) HandlerOf[any] {
		return HandlerFuncOf[any](func(ctx context.Context, e any) error {
			ctx = NewProducerContext(ctx, s.Producer)
			return next.Process(ctx, e)
		})
	}
}
