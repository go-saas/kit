package event

import (
	"context"
	"fmt"
	klog "github.com/go-kratos/kratos/v2/log"
	kerrors "github.com/goxiaoy/go-saas-kit/pkg/errors"
	"github.com/goxiaoy/uow"
)

var (
	_ Event = (*Message)(nil)
)

type Event interface {
	Key() string
	Value() []byte
}

type Message struct {
	key   string
	value []byte
}

func (m *Message) Key() string {
	return m.key
}

func (m *Message) Value() []byte {
	return m.value
}

func NewMessage(key string, value []byte) Event {
	return &Message{
		key:   key,
		value: value,
	}
}

type Handler func(context.Context, Event) error

type HandlerOf[T any] func(context.Context, T) error

type TransformerOf[T any] func(e Event) (T, error)

type Sender interface {
	Send(ctx context.Context, msg Event) error
	Close() error
}

type Receiver interface {
	Receive(ctx context.Context, handler Handler) error
	Close() error
}

//ChainHandler combine multiple handlers like http middleware
func ChainHandler(h ...Handler) Handler {
	return func(ctx context.Context, event Event) error {
		for _, handler := range h {
			if err := handler(ctx, event); err != nil {
				return err
			}
		}
		return nil
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

//RecoverHandler wrap next with recover. prevent consumer from panic
func RecoverHandler(next Handler, opt ...RecoverOption) Handler {
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
	return func(ctx context.Context, event Event) (err error) {
		defer func() {
			if rerr := recover(); rerr != nil {
				stack := kerrors.Stack(0)
				logger.Errorf("panic recovered: %s\n %s", rerr, stack)
				if rrerr, ok := rerr.(error); ok {
					err = op.formatter(ctx, fmt.Errorf("panic recovered: %w\n %s", rrerr, stack))
				} else {
					err = fmt.Errorf("panic recovered: %s\n %s", rerr, stack)
					err = op.formatter(ctx, err)
				}
			}
		}()
		err = next(ctx, event)
		if err == nil {
			return nil
		}
		return op.formatter(ctx, err)
	}
}

//UowHandler wrap handler into a unit of work (transaction)
func UowHandler(uowMgr uow.Manager, handler Handler) Handler {
	return func(ctx context.Context, event Event) error {
		return uowMgr.WithNew(ctx, func(ctx context.Context) error {
			return handler(ctx, event)
		})
	}
}

//FilterKeyHandler filter event by key compare
func FilterKeyHandler(key string, handler Handler) Handler {
	return func(ctx context.Context, event Event) error {
		if event.Key() == key {
			return handler(ctx, event)
		}
		return nil
	}
}

//TransformHandler transform Event into type generic T
func TransformHandler[T any](transformer TransformerOf[T], next HandlerOf[T]) Handler {
	return func(ctx context.Context, event Event) error {
		if data, err := transformer(event); err != nil {
			return err
		} else {
			return next(ctx, data)
		}
	}
}
