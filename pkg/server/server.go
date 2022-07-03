package server

import (
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-saas/kit/pkg/conf"
	kregistry "github.com/go-saas/kit/pkg/registry"
	"github.com/google/wire"
)

var DefaultProviderSet = wire.NewSet(wire.Value(ReqDecode), wire.Value(ResEncoder), wire.Value(ErrEncoder), NewRegistrar)

func NewRegistrar(services *conf.Services) (registry.Registrar, error) {
	r, _, err := kregistry.NewRegister(services.Registry)
	return r, err
}

var (
	ReqDecode  http.DecodeRequestFunc  = http.DefaultRequestDecoder
	ResEncoder http.EncodeResponseFunc = http.DefaultResponseEncoder
	ErrEncoder http.EncodeErrorFunc    = http.DefaultErrorEncoder
)
