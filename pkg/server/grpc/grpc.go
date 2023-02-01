package grpc

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-saas/kit/pkg/conf"
	"github.com/go-saas/kit/pkg/server/endpoint"
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
func PatchOpts(l log.Logger, opts []grpc.ServerOption, server *conf.Server) []grpc.ServerOption {

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

type Server struct {
	*grpc.Server
	cfg *conf.Server
}

func NewServer(cfg *conf.Server, opts ...grpc.ServerOption) *Server {
	return &Server{
		Server: grpc.NewServer(opts...),
		cfg:    cfg,
	}
}

func (s *Server) Endpoint() (url *url.URL, err error) {
	url, err = s.Server.Endpoint()
	if err != nil || url == nil {
		return
	}
	if s.cfg.Grpc != nil && len(s.cfg.Grpc.Endpoint) > 0 {
		//TODO tls
		return endpoint.NewEndpoint(endpoint.Scheme("grpc", false), s.cfg.Grpc.Endpoint), nil
	}
	return
}
