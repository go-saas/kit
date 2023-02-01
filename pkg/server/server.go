package server

import (
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-saas/kit/pkg/conf"
	kitdi "github.com/go-saas/kit/pkg/di"
	kregistry "github.com/go-saas/kit/pkg/registry"
	"github.com/goava/di"
	"google.golang.org/protobuf/types/known/durationpb"
	"time"
)

var DefaultProviderSet = kitdi.NewSet(
	kitdi.Value(ReqDecode),
	kitdi.Value(ResEncoder),
	kitdi.Value(ErrEncoder),
	NewRegistrar,
	NewWebMultiTenancyOption,
)

const (
	defaultSrvName = "default"
)

var (
	defaultServerConfig = &conf.Server{
		Http: &conf.Server_HTTP{
			Addr:    ":9080",
			Timeout: durationpb.New(5 * time.Second),
		},
		Grpc: &conf.Server_GRPC{
			Addr:    ":9081",
			Timeout: durationpb.New(5 * time.Second),
		},
	}
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
