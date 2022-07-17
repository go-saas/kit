package server

import (
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-saas/kit/pkg/conf"
	kitdi "github.com/go-saas/kit/pkg/di"
	kregistry "github.com/go-saas/kit/pkg/registry"
	"github.com/goava/di"
)

var DefaultProviderSet = kitdi.NewSet(
	kitdi.Value(ReqDecode),
	kitdi.Value(ResEncoder),
	kitdi.Value(ErrEncoder),
	NewRegistrar,
)

func NewRegistrar(services *conf.Services, container *di.Container) (registry.Registrar, error) {
	err := container.Provide(func() *kregistry.Config { return services.Registry })
	if err != nil {
		return nil, err
	}
	r, _, err := kregistry.NewRegister(services.Registry, container)
	return r, err
}

var (
	ReqDecode  http.DecodeRequestFunc  = http.DefaultRequestDecoder
	ResEncoder http.EncodeResponseFunc = http.DefaultResponseEncoder
	ErrEncoder http.EncodeErrorFunc    = http.DefaultErrorEncoder
)
