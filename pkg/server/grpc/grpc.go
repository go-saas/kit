package grpc

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-saas/kit/pkg/conf"
	"github.com/go-saas/kit/pkg/server/common"
	"google.golang.org/protobuf/proto"
	"net/url"
)

type (
	// ServiceRegister register grpc handler into grpc server
	ServiceRegister interface {
		Register(server *grpc.Server, middleware ...middleware.Middleware)
	}
	ServiceRegisterFunc func(server *grpc.Server, middleware ...middleware.Middleware)
)

func (f ServiceRegisterFunc) Register(server *grpc.Server, middleware ...middleware.Middleware) {
	f(server, middleware...)
}

func ChainServiceRegister(r ...ServiceRegister) ServiceRegister {
	return ServiceRegisterFunc(func(server *grpc.Server, middleware ...middleware.Middleware) {
		for _, register := range r {
			register.Register(server, middleware...)
		}
	})
}

// PatchOpts Patch grpc options with given service name and configs
func PatchOpts(l log.Logger, opts []grpc.ServerOption, name string, services *conf.Services) []grpc.ServerOption {
	//default config
	server := proto.Clone(common.DefaultServerConfig).(*conf.Server)
	if def, ok := services.Servers[common.DefaultSrvName]; ok {
		//merge default config
		proto.Merge(server, def)
	}
	if s, ok := services.Servers[name]; ok {
		//merge service config
		proto.Merge(server, s)
	}

	if server.Grpc.Network != "" {
		opts = append(opts, grpc.Network(server.Grpc.Network))
	}
	if server.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(server.Grpc.Addr))
	}
	if server.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(server.Grpc.Timeout.AsDuration()))
	}
	return opts
}

type server struct {
	*grpc.Server
	cfg *conf.Dev
}

func NewServer(cfg *conf.Dev, opts ...grpc.ServerOption) *server {
	return &server{
		Server: grpc.NewServer(opts...),
		cfg:    nil,
	}
}

func (s *server) Endpoint() (*url.URL, error) {
	url, err := s.Server.Endpoint()
	if err != nil || url == nil || !s.cfg.Docker {
		return url, err
	}
	//replace host
	url.Host = "host.docker.internal"
	return url, err
}
