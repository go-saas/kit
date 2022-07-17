package event

import (
	"context"
	"fmt"
	"github.com/goava/di"
	"sync"
)

type LazyConsumer func(ctx context.Context, c *Config, container *di.Container) (Consumer, error)
type LazyProducer func(c *Config, container *di.Container) (*ProducerMux, error)

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

func resolveType(cfg *Config) (string, error) {
	t := cfg.Type
	if len(t) > 0 {
		return t, nil
	}
	if cfg.Kafka != nil {
		t = "kafka"
	} else if cfg.Pulsar != nil {
		t = "pulsar"
	}
	if len(t) == 0 {
		return t, fmt.Errorf("can not reolve event type %s", t)
	}
	return t, nil
}
