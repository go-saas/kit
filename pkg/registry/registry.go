package registry

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/registry"
	"sync"
)

type Factory func(c *Config) (registry.Registrar, registry.Discovery, error)

var (
	_registryMap  = map[string]Factory{}
	_registryLock = sync.RWMutex{}
)

func Register(kind string, factory Factory) {
	_registryLock.Lock()
	defer _registryLock.Unlock()
	_registryMap[kind] = factory
}

func NewRegister(c *Config) (registry.Registrar, registry.Discovery, error) {
	_registryLock.RLock()
	defer _registryLock.RUnlock()
	if len(c.Type) == 0 {
		return nil, nil, fmt.Errorf("registry type is required")
	}
	r, ok := _registryMap[c.Type]
	if !ok {
		return nil, nil, fmt.Errorf("registry type %s not found", c.Type)
	}
	return r(c)
}
