package event

import (
	"context"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	"io"
	"sync"
)

type Consumer interface {
	io.Closer
	Process(ctx context.Context, handler ConsumerHandler) error
}

type ConsumerHandler HandlerOf[Event]

type ConsumerHandlerFunc HandlerFuncOf[Event]

func (h ConsumerHandlerFunc) Process(ctx context.Context, e Event) error {
	return h(ctx, e)
}

// NewTransformer wrap handle by transform event to T
func NewTransformer[T any](t func(context.Context, Event) (T, error), f HandlerOf[T]) ConsumerHandler {
	return ConsumerHandlerFunc(func(ctx context.Context, e Event) error {
		tt, err := t(ctx, e)
		if err != nil {
			return err
		}
		return f.Process(ctx, tt)
	})
}

type ConsumerMiddlewareFunc func(ConsumerHandler) ConsumerHandler

func ConsumerChain(m ...ConsumerMiddlewareFunc) ConsumerMiddlewareFunc {
	return func(next ConsumerHandler) ConsumerHandler {
		for i := len(m) - 1; i >= 0; i-- {
			next = m[i](next)
		}
		return next
	}
}

func FilterKey(key string, handler ConsumerHandler) ConsumerHandler {
	return ConsumerHandlerFunc(func(ctx context.Context, event Event) error {
		if event.Key() == key {
			return handler.Process(ctx, event)
		}
		return nil
	})
}

type ServeMux struct {
	mu      sync.RWMutex
	mws     []ConsumerMiddlewareFunc
	handles []ConsumerHandler
	r       Consumer
}

// Use appends a ConsumerMiddlewareFunc to the chain.
// Middlewares are executed in the order that they are applied to the ServeMux.
func (mux *ServeMux) Use(mws ...ConsumerMiddlewareFunc) {
	mux.mu.Lock()
	defer mux.mu.Unlock()
	for _, fn := range mws {
		mux.mws = append(mux.mws, fn)
	}
}

// Append will append handler into mux,
func (mux *ServeMux) Append(h ConsumerHandler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()
	mux.handles = append(mux.handles, h)
}

// Process call handler one by one until error happens
func (mux *ServeMux) Process(ctx context.Context, event Event) error {
	mux.mu.RLock()
	defer mux.mu.RUnlock()

	//push consumer into context
	ctx = NewConsumerContext(ctx, mux.r)

	c := ConsumerChain(mux.mws...)
	h := c(ConsumerHandlerFunc(func(ctx context.Context, e Event) error {
		for _, handle := range mux.handles {
			if err := handle.Process(ctx, event); err != nil {
				return err
			}
		}
		return nil
	}))
	return h.Process(ctx, event)

}

// ConsumerFactoryServer resolve LazyConsumer from factory, then wrap as kratos server
type ConsumerFactoryServer struct {
	*ServeMux
	lr  LazyConsumer
	cfg *conf.Event
}

var _ transport.Server = (*ConsumerFactoryServer)(nil)

func NewConsumerFactoryServer(cfg *conf.Event) *ConsumerFactoryServer {
	_typeConsumerMux.RLock()
	defer _typeConsumerMux.RUnlock()
	var r LazyConsumer
	var ok bool
	if r, ok = _typeConsumerRegister[cfg.Type]; !ok {
		panic(cfg.Type + " event server not registered")
	}
	return &ConsumerFactoryServer{
		ServeMux: &ServeMux{},
		cfg:      cfg,
		lr:       r,
	}
}

func (f *ConsumerFactoryServer) Start(ctx context.Context) error {
	if f.r != nil {
		panic("server can not start twice")
	}
	r, err := f.lr(ctx, f.cfg)
	if err != nil {
		return err
	}
	f.r = r
	return r.Process(ctx, f)
}

func (f *ConsumerFactoryServer) Stop(ctx context.Context) error {
	return f.r.Close()
}

type ConsumerServer struct {
	*ServeMux
}

var _ transport.Server = (*ConsumerServer)(nil)

// NewConsumerServer create server from Consumer directly
func NewConsumerServer(r Consumer) *ConsumerServer {
	return &ConsumerServer{&ServeMux{r: r}}
}

func (s *ConsumerServer) Start(ctx context.Context) error {
	return s.r.Process(ctx, s)
}

func (s *ConsumerServer) Stop(ctx context.Context) error {
	return s.r.Close()
}
