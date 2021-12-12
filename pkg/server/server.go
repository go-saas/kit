package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/handlers"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	"net"
	"strings"
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
		allowMethods := []string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE", "PATCH"}
		allowMethods = append(allowMethods, server.Http.Cors.GetAllowedMethods()...)
		opts = append(opts, http.Filter(handlers.CORS(
			handlers.AllowedOrigins(server.Http.Cors.GetAllowedOrigins()),
			handlers.AllowedMethods(allowMethods),
			handlers.AllowedHeaders(append([]string{"Content-Type", "Authorization"}, server.Http.Cors.AllowedHeaders...)),
		)))
	}
	return opts
}

func ClientIP(ctx context.Context) string {
	if t, ok := transport.FromServerContext(ctx); ok {
		if ht, ok := t.(*http.Transport); ok {
			xForwardedFor := ht.Request().Header.Get("X-Forwarded-For")
			ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
			if ip != "" {
				return ip
			}

			ip = strings.TrimSpace(ht.Request().Header.Get("X-Real-Ip"))
			if ip != "" {
				return ip
			}

			if ip, _, err := net.SplitHostPort(strings.TrimSpace(ht.Request().RemoteAddr)); err == nil {
				return ip
			}
		}
	}
	return ""
}

func ClientUserAgent(ctx context.Context) string {
	if t, ok := transport.FromServerContext(ctx); ok {
		if ht, ok := t.(*http.Transport); ok {
			return ht.Request().UserAgent()
		}
	}
	return ""
}
