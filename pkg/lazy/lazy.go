package lazy

import (
	"context"
	"sync"
)

type Of[T any] struct {
	factory  func(ctx context.Context) (T, error)
	lockItem sync.Mutex
	value    T
}

func (l *Of[T]) Value(ctx context.Context) (T, error) {
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

func New[T any](factory func(ctx context.Context) (T, error)) *Of[T] {
	return &Of[T]{factory: factory}
}
