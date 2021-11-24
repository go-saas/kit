package server

import (
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/handlers"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
)

// PatchGrpcOpts Patch grpc options with given service name and configs
func PatchGrpcOpts(opts []grpc.ServerOption, name string, services *conf.Services) []grpc.ServerOption {
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

// PatchHttpOpts Patch http options with given service name and configs
func PatchHttpOpts(opts []http.ServerOption, name string, services *conf.Services) []http.ServerOption {
	server, ok := services.Servers[name]
	if !ok {
		panic(errors.New(fmt.Sprintf(" %v server not found", name)))
	}
	if server.Http.Network != "" {
		opts = append(opts, http.Network(server.Http.Network))
	}
	if server.Http.Addr != "" {
		opts = append(opts, http.Address(server.Http.Addr))
	}
	if server.Http.Timeout != nil {
		opts = append(opts, http.Timeout(server.Http.Timeout.AsDuration()))
	}
	if server.Http.Cors != nil {
		opts = append(opts, http.Filter(handlers.CORS(
			handlers.AllowedOrigins(server.Http.Cors.GetAllowedOrigins()),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE", "PATCH"}),
			handlers.AllowedMethods(server.Http.Cors.GetAllowedMethods()),
			handlers.AllowedHeaders(append([]string{"Content-Type", "Authorization"}, server.Http.Cors.AllowedHeaders...)),
		)))
	}
	return opts
}
