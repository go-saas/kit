package event

import (
	"context"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/goava/di"
	"io"
	"sync"
)

type Consumer interface {
	io.Closer
	// Process start process event with handler
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

type ConsumerMux struct {
	mu      sync.RWMutex
	mws     []ConsumerMiddlewareFunc
	handles []ConsumerHandler
	r       Consumer
}

// Use appends a ConsumerMiddlewareFunc to the chain.
// Middlewares are executed in the order that they are applied to the ConsumerMux.
func (mux *ConsumerMux) Use(mws ...ConsumerMiddlewareFunc) {
	mux.mu.Lock()
	defer mux.mu.Unlock()
	for _, fn := range mws {
		mux.mws = append(mux.mws, fn)
	}
}

// Append will append handler into mux,
func (mux *ConsumerMux) Append(h ConsumerHandler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()
	mux.handles = append(mux.handles, h)
}

// Process call handler one by one until error happens
func (mux *ConsumerMux) Process(ctx context.Context, event Event) error {
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
	*ConsumerMux
	lr        LazyConsumer
	cfg       *Config
	container *di.Container
}

var _ transport.Server = (*ConsumerFactoryServer)(nil)

func NewConsumerFactoryServer(cfg *Config, container *di.Container) *ConsumerFactoryServer {
	_typeConsumerMux.RLock()
	defer _typeConsumerMux.RUnlock()
	var r LazyConsumer
	var ok bool
	t, err := resolveType(cfg)
	if err != nil {
		panic(err)
	}

	if r, ok = _typeConsumerRegister[t]; !ok {
		panic(cfg.Type + " event server not registered")
	}
	return &ConsumerFactoryServer{
		ConsumerMux: &ConsumerMux{},
		cfg:         cfg,
		lr:          r,
		container:   container,
	}
}

func (f *ConsumerFactoryServer) Start(ctx context.Context) error {
	if f.r != nil {
		panic("server can not start twice")
	}
	r, err := f.lr(ctx, f.cfg, f.container)
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
	*ConsumerMux
}

var _ transport.Server = (*ConsumerServer)(nil)

// NewConsumerServer create server from Consumer directly
func NewConsumerServer(r Consumer) *ConsumerServer {
	return &ConsumerServer{&ConsumerMux{r: r}}
}

func (s *ConsumerServer) Start(ctx context.Context) error {
	return s.r.Process(ctx, s)
}

func (s *ConsumerServer) Stop(ctx context.Context) error {
	return s.r.Close()
}
