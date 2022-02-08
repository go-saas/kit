package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/handlers"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	"github.com/goxiaoy/go-saas-kit/pkg/csrf"
	"net"
	"net/http"
	"strings"
)

func ResolveHttpRequest(ctx context.Context) (*http.Request, bool) {
	if t, ok := transport.FromServerContext(ctx); ok {
		if ht, ok := t.(*khttp.Transport); ok {
			return ht.Request(), true
		}
	}
	return nil, false
}

// PatchHttpOpts Patch khttp options with given service name and configs. f use global filters
func PatchHttpOpts(l log.Logger,
	opts []khttp.ServerOption,
	name string,
	services *conf.Services,
	sCfg *conf.Security,
	reqDecoder khttp.DecodeRequestFunc,
	resEncoder khttp.EncodeResponseFunc,
	errEncoder khttp.EncodeErrorFunc,
	f ...khttp.FilterFunc) []khttp.ServerOption {
	server, ok := services.Servers[name]
	if !ok {
		panic(errors.New(fmt.Sprintf(" %v server not found", name)))
	}
	if server.Http.Network != "" {
		opts = append(opts, khttp.Network(server.Http.Network))
	}
	if server.Http.Addr != "" {
		opts = append(opts, khttp.Address(server.Http.Addr))
	}
	if server.Http.Timeout != nil {
		opts = append(opts, khttp.Timeout(server.Http.Timeout.AsDuration()))
	}
	if reqDecoder != nil {
		opts = append(opts, khttp.RequestDecoder(reqDecoder))
	}
	if resEncoder != nil {
		opts = append(opts, khttp.ResponseEncoder(resEncoder))
	}
	if errEncoder != nil {
		opts = append(opts, khttp.ErrorEncoder(errEncoder))
	}
	var filters []khttp.FilterFunc

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
	opts = append(opts, khttp.Filter(filters...))
	return opts
}

func ClientIP(ctx context.Context) string {
	if r, ok := ResolveHttpRequest(ctx); ok {
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
	if r, ok := ResolveHttpRequest(ctx); ok {
		return r.UserAgent()
	}
	return ""
}

func IsSecure(ctx context.Context) bool {
	if r, ok := ResolveHttpRequest(ctx); ok {
		return r.URL.Scheme == "https"
	}
	return false
}

func IsWebsocket(ctx context.Context) bool {
	if r, ok := ResolveHttpRequest(ctx); ok {
		h := r.Header["Upgrade"]
		return len(h) > 0 && h[0] == "websocket"
	}
	return false
}

func IsAjax(ctx context.Context) bool {
	if r, ok := ResolveHttpRequest(ctx); ok {
		h := r.Header["X-Requested-With"]
		return len(h) > 0 && h[0] == "XMLHttpRequest"
	}
	return false
}

func SetCookie(ctx context.Context, value string) bool {
	if t, ok := transport.FromServerContext(ctx); ok {
		if ht, ok := t.(*khttp.Transport); ok {
			ht.ReplyHeader().Set("Set-Cookie", value)
			return true
		}
	}
	return false
}

type ErrorHandler interface {
	Wrap(func(w http.ResponseWriter, r *http.Request) error) http.Handler
}

type DefaultErrorHandler struct {
	errEncoder khttp.EncodeErrorFunc
}

var _ ErrorHandler = (*DefaultErrorHandler)(nil)

func NewDefaultErrorHandler(errEncoder khttp.EncodeErrorFunc) *DefaultErrorHandler {
	return &DefaultErrorHandler{errEncoder: errEncoder}
}

func (e *DefaultErrorHandler) Wrap(f func(w http.ResponseWriter, r *http.Request) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			e.errEncoder(w, r, err)
			return
		}
	})
}
