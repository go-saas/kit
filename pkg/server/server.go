package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/handlers"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	"github.com/goxiaoy/go-saas-kit/pkg/csrf"
	"github.com/goxiaoy/go-saas-kit/pkg/kratos"
	"net"
	"strings"
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

// PatchHttpOpts Patch http options with given service name and configs. f use global filters
func PatchHttpOpts(l log.Logger,
	opts []http.ServerOption,
	name string,
	services *conf.Services,
	sCfg *conf.Security,
	reqDecoder http.DecodeRequestFunc,
	resEncoder http.EncodeResponseFunc,
	errEncoder http.EncodeErrorFunc,
	f ...http.FilterFunc) []http.ServerOption {
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
	if reqDecoder != nil {
		opts = append(opts, http.RequestDecoder(reqDecoder))
	}
	if resEncoder != nil {
		opts = append(opts, http.ResponseEncoder(resEncoder))
	}
	if errEncoder != nil {
		opts = append(opts, http.ErrorEncoder(errEncoder))
	}
	var filters []http.FilterFunc

	if server.Http.Cors != nil {
		allowMethods := []string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE", "PATCH"}
		allowMethods = append(allowMethods, server.Http.Cors.GetAllowedMethods()...)
		filters = append(filters, handlers.CORS(
			handlers.AllowedOrigins(server.Http.Cors.GetAllowedOrigins()),
			handlers.AllowedMethods(allowMethods),
			handlers.AllowedHeaders(append([]string{"Content-Type", "Authorization"}, server.Http.Cors.AllowedHeaders...)),
		))
	}
	if server.Http.Csrf != nil {
		filters = append(filters, csrf.NewCsrf(l, sCfg, server.Http.Csrf, errEncoder))
	}
	filters = append(filters, f...)
	opts = append(opts, http.Filter(filters...))
	return opts
}

func ClientIP(ctx context.Context) string {
	if r, ok := kratos.ResolveHttpRequest(ctx); ok {
		xForwardedFor := r.Header.Get("X-Forwarded-For")
		ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
		if ip != "" {
			return ip
		}

		ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
		if ip != "" {
			return ip
		}

		if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
			return ip
		}
	}
	return ""
}

func ClientUserAgent(ctx context.Context) string {
	if r, ok := kratos.ResolveHttpRequest(ctx); ok {
		return r.UserAgent()
	}
	return ""
}
