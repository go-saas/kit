package server

import (
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-saas/kit/pkg/conf"
	kitdi "github.com/go-saas/kit/pkg/di"
	kregistry "github.com/go-saas/kit/pkg/registry"
)

var DefaultProviderSet = kitdi.NewSet(
	func() http.DecodeRequestFunc { return ReqDecode },
	func() http.EncodeResponseFunc { return ResEncoder },
	func() http.EncodeErrorFunc { return ErrEncoder },
	NewRegistrar,
)

func NewRegistrar(services *conf.Services) (registry.Registrar, error) {
	r, _, err := kregistry.NewRegister(services.Registry)
	return r, err
}

var (
	ReqDecode  http.DecodeRequestFunc  = http.DefaultRequestDecoder
	ResEncoder http.EncodeResponseFunc = http.DefaultResponseEncoder
	ErrEncoder http.EncodeErrorFunc    = http.DefaultErrorEncoder
)
