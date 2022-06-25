package event

import (
	"context"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	"sync"
)

type LazyConsumer func(ctx context.Context, c *conf.Event) (Consumer, error)
type LazyProducer func(c *conf.Event) (*ProducerMux, error)

var (
	_typeConsumerMux      sync.RWMutex
	_typeConsumerRegister map[string]LazyConsumer

	_typeProducerMux      sync.RWMutex
	_typeProducerRegister map[string]LazyProducer
)

func init() {
	_typeConsumerRegister = map[string]LazyConsumer{}
	_typeProducerRegister = map[string]LazyProducer{}
}

func RegisterConsumer(kind string, e LazyConsumer) {
	_typeConsumerMux.Lock()
	defer _typeConsumerMux.Unlock()
	if len(kind) == 0 {
		panic("kind is required")
	}
	_typeConsumerRegister[kind] = e
}

func RegisterProducer(kind string, e LazyProducer) {
	_typeProducerMux.Lock()
	defer _typeProducerMux.Unlock()
	if len(kind) == 0 {
		panic("kind is required")
	}
	_typeProducerRegister[kind] = e
}
