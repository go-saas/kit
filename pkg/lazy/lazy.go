package lazy

import (
	"context"
	"sync"
)

type Of[T any] interface {
	Value(ctx context.Context) (T, error)
	Initialized() bool
}

type of[T any] struct {
	factory  func(ctx context.Context) (T, error)
	lockItem sync.Mutex
	value    T
}

func (l *of[T]) Value(ctx context.Context) (T, error) {
	l.lockItem.Lock()
	defer l.lockItem.Unlock()
	var err error
	if l.factory != nil {
		l.value, err = l.factory(ctx)
		if err != nil {
			l.factory = nil
		} else {
			return l.value, err
		}
	}
	return l.value, err
}

func (l *of[T]) Initialized() bool {
	l.lockItem.Lock()
	defer l.lockItem.Unlock()
	return l.factory == nil
}

func New[T any](factory func(ctx context.Context) (T, error)) Of[T] {
	if factory == nil {
		panic("factory can not be nil")
	}
	return &of[T]{factory: factory}
}
