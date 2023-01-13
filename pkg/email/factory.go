package email

import (
	"github.com/goava/di"
	"sync"
)

type ProviderFunc func(c *Config, container *di.Container) (Client, error)

var (
	_typeMux      sync.RWMutex
	_typeRegister = map[string]ProviderFunc{}
)

func RegisterProvider(kind string, f ProviderFunc) {
	if len(kind) == 0 {
		panic("kind is required")
	}
	_typeMux.Lock()
	defer _typeMux.Unlock()
	_typeRegister[kind] = f
}
