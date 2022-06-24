package event

import (
	"context"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	"sync"
)

type LazyReceiver func(ctx context.Context, c *conf.Event) (Receiver, error)
type LazySender func(c *conf.Event) (Sender, func(), error)

var (
	_typeReceiverMux sync.RWMutex
	_typeRegister    map[string]LazyReceiver

	_typeSenderMux      sync.RWMutex
	_typeSenderRegister map[string]LazySender
)

func init() {
	_typeRegister = map[string]LazyReceiver{}
	_typeSenderRegister = map[string]LazySender{}
}

func RegisterReceiver(kind string, e LazyReceiver) {
	_typeReceiverMux.Lock()
	defer _typeReceiverMux.Unlock()
	if len(kind) == 0 {
		panic("kind is required")
	}
	_typeRegister[kind] = e
}

func RegisterSender(kind string, e LazySender) {
	_typeSenderMux.Lock()
	defer _typeSenderMux.Unlock()
	if len(kind) == 0 {
		panic("kind is required")
	}
	_typeSenderRegister[kind] = e
}

func NewFactoryServer(cfg *conf.Event) *FactoryServer {
	_typeReceiverMux.RLock()
	defer _typeReceiverMux.RUnlock()
	var r LazyReceiver
	var ok bool
	if r, ok = _typeRegister[cfg.Type]; !ok {
		panic(cfg.Type + " event server not registered")
	}
	return &FactoryServer{
		ServeMux: &ServeMux{},
		cfg:      cfg,
		lr:       r,
	}
}
