package event

import (
	"context"
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

type Sender interface {
	Send(ctx context.Context, msg Event) error
	Close() error
}

type Receiver interface {
	Receive(ctx context.Context, handler Handler) error
	Close() error
}

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

//UowHandler wrap handler into a unit of work
func UowHandler(uowMgr uow.Manager, handler Handler) Handler {
	return func(ctx context.Context, event Event) error {
		return uowMgr.WithNew(ctx, func(ctx context.Context) error {
			return handler(ctx, event)
		})
	}
}

func FilterKeyHandler(key string, handler Handler) Handler {
	return func(ctx context.Context, event Event) error {
		if event.Key() == key {
			return handler(ctx, event)
		}
		return nil
	}
}

//TransformHandler transform Event into some type
func TransformHandler[T any](transformer func(e Event) (T, error), next func(context.Context, T) error) Handler {
	return func(ctx context.Context, event Event) error {
		if data, err := transformer(event); err != nil {
			return err
		} else {
			return next(ctx, data)
		}
	}
}
