package registry

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/goava/di"
	"sync"
)

type Factory func(c *Config, container *di.Container) (registry.Registrar, registry.Discovery, error)

var (
	_registryMap  = map[string]Factory{}
	_registryLock = sync.RWMutex{}
)

func Register(kind string, factory Factory) {
	_registryLock.Lock()
	defer _registryLock.Unlock()
	_registryMap[kind] = factory
}

func NewRegister(c *Config, container *di.Container) (registry.Registrar, registry.Discovery, error) {
	_registryLock.RLock()
	defer _registryLock.RUnlock()
	if len(c.Type) == 0 {
		return nil, nil, fmt.Errorf("registry type is required")
	}
	r, ok := _registryMap[c.Type]
	if !ok {
		return nil, nil, fmt.Errorf("registry type %s not found", c.Type)
	}
	return r(c, container)
}

// Discovery is service discovery.
type Discovery interface {
	registry.Discovery
	// WatchAll creates a watcher to all services
	WatchAll(ctx context.Context) (registry.Watcher, error)
}
