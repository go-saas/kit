package server

import (
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
)

// PatchGrpcOpts Patch grpc options with given service name and configs
func PatchGrpcOpts(l log.Logger, opts []grpc.ServerOption, name string, services *conf.Services) []grpc.ServerOption {
	server, ok := services.Servers[name]
	if !ok {
		panic(errors.New(fmt.Sprintf(" %v server not found", name)))
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
