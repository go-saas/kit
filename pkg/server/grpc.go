package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-saas/kit/pkg/conf"
	"google.golang.org/protobuf/proto"
)

type (
	// GrpcServiceRegister register grpc handler into grpc server
	GrpcServiceRegister interface {
		Register(server *grpc.Server, middleware middleware.Middleware)
	}
	GrpcServiceRegisterFunc func(server *grpc.Server, middleware middleware.Middleware)
)

func (f GrpcServiceRegisterFunc) Register(server *grpc.Server, middleware middleware.Middleware) {
	f(server, middleware)
}

func ChainGrpcServiceRegister(r ...GrpcServiceRegister) GrpcServiceRegister {
	return GrpcServiceRegisterFunc(func(server *grpc.Server, middleware middleware.Middleware) {
		for _, register := range r {
			register.Register(server, middleware)
		}
	})
}

// PatchGrpcOpts Patch grpc options with given service name and configs
func PatchGrpcOpts(l log.Logger, opts []grpc.ServerOption, name string, services *conf.Services) []grpc.ServerOption {
	//default config
	server := proto.Clone(defaultServiceConfig).(*conf.Server)
	if def, ok := services.Servers[defaultSrvName]; ok {
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
