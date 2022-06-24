package event

import (
	"context"
	"fmt"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	"github.com/goxiaoy/uow"
	"io"
	"net/http"
	"sync"
)

var (
	_ Event = (*Message)(nil)
)

type Header interface {
	Get(key string) string
	Set(key string, value string)
	Keys() []string
}

type headerCarrier http.Header

// Get returns the value associated with the passed key.
func (hc headerCarrier) Get(key string) string {
	return http.Header(hc).Get(key)
}

// Set stores the key-value pair.
func (hc headerCarrier) Set(key string, value string) {
	http.Header(hc).Set(key, value)
}

// Keys lists the keys stored in this carrier.
func (hc headerCarrier) Keys() []string {
	keys := make([]string, 0, len(hc))
	for k := range http.Header(hc) {
		keys = append(keys, k)
	}
	return keys
}

type Event interface {
	Header() Header
	Key() string
	Value() []byte
}

type Message struct {
	header headerCarrier
	key    string
	value  []byte
}

func (m *Message) Key() string {
	return m.key
}

func (m *Message) Header() Header {
	return m.header
}

func (m *Message) Value() []byte {
	return m.value
}

func NewMessage(key string, value []byte) Event {
	return &Message{
		key:    key,
		value:  value,
		header: headerCarrier{},
	}
}

type Receiver interface {
	io.Closer
	Receive(ctx context.Context, handler Handler) error
}

type HandlerOf[T any] interface {
	Process(context.Context, T) error
}
type Handler HandlerOf[Event]

type HandlerFuncOf[T any] func(context.Context, T) error

func (h HandlerFuncOf[T]) Process(ctx context.Context, e T) error {
	return h(ctx, e)
}

type HandlerFunc HandlerFuncOf[Event]

func (h HandlerFunc) Process(ctx context.Context, e Event) error {
	return h(ctx, e)
}

// NewTransformer wrap handle by transform event to T
func NewTransformer[T any](t func(context.Context, Event) (T, error), f HandlerOf[T]) Handler {
	return HandlerFunc(func(ctx context.Context, e Event) error {
		tt, err := t(ctx, e)
		if err != nil {
			return err
		}
		return f.Process(ctx, tt)
	})
}

type MiddlewareFunc func(Handler) Handler

func Chain(m ...MiddlewareFunc) MiddlewareFunc {
	return func(next Handler) Handler {
		for i := len(m) - 1; i >= 0; i-- {
			next = m[i](next)
		}
		return next
	}
}

type RecoverOption func(*recoverOptions)

type recoverOptions struct {
	formatter ErrFormatFunc
	logger    klog.Logger
}

type ErrFormatFunc func(ctx context.Context, err error) error

func WithErrorFormatter(f ErrFormatFunc) RecoverOption {
	return func(o *recoverOptions) {
		o.formatter = f
	}
}

func WithLogger(logger klog.Logger) RecoverOption {
	return func(o *recoverOptions) {
		o.logger = logger
	}
}

// Logging logging errors
func Logging(logger klog.Logger) MiddlewareFunc {
	return func(next Handler) Handler {
		return HandlerFunc(func(ctx context.Context, event Event) error {
			err := next.Process(ctx, event)
			if err != nil {
				_ = klog.WithContext(ctx, logger).Log(klog.LevelError,
					klog.DefaultMessageKey, err.Error(),
					"event", event.Key())
			} else {
				_ = klog.WithContext(ctx, logger).Log(klog.LevelInfo,
					"event", event.Key())
			}
			return err
		})
	}
}

//Recover prevent consumer from panic
func Recover(opt ...RecoverOption) MiddlewareFunc {
	op := recoverOptions{
		logger: klog.GetLogger(),
		formatter: func(ctx context.Context, err error) error {
			return err
		},
	}
	for _, o := range opt {
		o(&op)
	}
	logger := klog.NewHelper(op.logger)
	return func(next Handler) Handler {
		return HandlerFunc(func(ctx context.Context, event Event) (err error) {
			defer func() {
				if rerr := recover(); rerr != nil {
					if rrerr, ok := rerr.(error); ok {
						wrrer := fmt.Errorf("panic recovered: %w", rrerr)
						logger.Error(wrrer)
						err = op.formatter(ctx, wrrer)
					} else {
						err = fmt.Errorf("panic recovered: %s", rerr)
						logger.Error(err)
						err = op.formatter(ctx, err)
					}
				}
			}()
			err = next.Process(ctx, event)
			if err == nil {
				return nil
			}
			return op.formatter(ctx, err)
		})
	}
}

//Uow wrap handler into a unit of work (transaction)
func Uow(uowMgr uow.Manager) MiddlewareFunc {
	return func(handler Handler) Handler {
		return HandlerFunc(func(ctx context.Context, event Event) error {
			return uowMgr.WithNew(ctx, func(ctx context.Context) error {
				return handler.Process(ctx, event)
			})
		})
	}
}

func FilterKey(key string, handler Handler) Handler {
	return HandlerFunc(func(ctx context.Context, event Event) error {
		if event.Key() == key {
			return handler.Process(ctx, event)
		}
		return nil
	})
}

type ServeMux struct {
	mu      sync.RWMutex
	mws     []MiddlewareFunc
	handles []Handler
	r       Receiver
}

// Use appends a MiddlewareFunc to the chain.
// Middlewares are executed in the order that they are applied to the ServeMux.
func (mux *ServeMux) Use(mws ...MiddlewareFunc) {
	mux.mu.Lock()
	defer mux.mu.Unlock()
	for _, fn := range mws {
		mux.mws = append(mux.mws, fn)
	}
}

// Append will append handler into mux,
func (mux *ServeMux) Append(h Handler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()
	mux.handles = append(mux.handles, h)
}

// Process call handler one by one until error happens
func (mux *ServeMux) Process(ctx context.Context, event Event) error {
	mux.mu.RLock()
	defer mux.mu.RUnlock()

	//push receiver into context
	ctx = NewReceiverContext(ctx, mux.r)

	c := Chain(mux.mws...)
	h := c(HandlerFunc(func(ctx context.Context, e Event) error {
		for _, handle := range mux.handles {
			if err := handle.Process(ctx, event); err != nil {
				return err
			}
		}
		return nil
	}))
	return h.Process(ctx, event)

}

func NewFactorySender(cfg *conf.Event) (Sender, func(), error) {
	_typeSenderMux.RLock()
	defer _typeSenderMux.RUnlock()
	if r, ok := _typeSenderRegister[cfg.Type]; !ok {
		panic(cfg.Type + " event sender not registered")
	} else {
		return r(cfg)
	}
}
