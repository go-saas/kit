package server

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authz"
	"github.com/goxiaoy/go-saas-kit/pkg/blob"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	"github.com/goxiaoy/go-saas-kit/pkg/csrf"
	"github.com/spf13/afero"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	defaultSrvName = "default"
)

var (
	defaultServiceConfig = &conf.Server{
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

type (
	// HttpServiceRegister register http handler into http server
	HttpServiceRegister interface {
		Register(server *khttp.Server, middleware middleware.Middleware)
	}
	HttpServiceRegisterFunc func(server *khttp.Server, middleware middleware.Middleware)
)

func (f HttpServiceRegisterFunc) Register(server *khttp.Server, middleware middleware.Middleware) {
	f(server, middleware)
}

func ChainHttpServiceRegister(r ...HttpServiceRegister) HttpServiceRegister {
	return HttpServiceRegisterFunc(func(server *khttp.Server, middleware middleware.Middleware) {
		for _, register := range r {
			register.Register(server, middleware)
		}
	})
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
			handlers.AllowCredentials(),
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

func HandleBlobs(basePath string, cfg blob.Config, srv *khttp.Server, factory blob.Factory) {
	if cfg == nil {
		return
	}
	router := mux.NewRouter()
	for s, config := range cfg {
		//local file
		handleBlob(s, config, factory, router)
	}
	if basePath == "" {
		basePath = "/assets"
	}
	if !strings.HasPrefix(basePath, "/") {
		basePath = fmt.Sprintf("/%s", basePath)
	}
	srv.HandlePrefix(basePath, router)
}
func handleBlob(name string, config *blob.BlobConfig, factory blob.Factory, router *mux.Router) {
	a := factory.Get(context.Background(), name, false).GetAfero()
	basePath := fmt.Sprintf("/%s", strings.TrimPrefix(config.BasePath, "/"))
	router.
		PathPrefix(basePath).
		Handler(http.StripPrefix(basePath, http.FileServer(http.FS(afero.NewIOFS(a.Fs)))))

}

func ResolveHttpRequest(ctx context.Context) (*http.Request, bool) {
	if t, ok := transport.FromServerContext(ctx); ok {
		if ht, ok := t.(*khttp.Transport); ok {
			return ht.Request(), true
		}
	}
	return nil, false
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

func Host(ctx context.Context) string {
	if r, ok := ResolveHttpRequest(ctx); ok {
		return r.Host
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

// AuthzGuardian guard http.Handler with authz
func AuthzGuardian(srv authz.Service, requirement authz.RequirementList, encoder khttp.EncodeErrorFunc, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		r, err := srv.BatchCheck(request.Context(), requirement)
		if err != nil {
			encoder(writer, request, err)
			return
		}
		for _, result := range r {
			if !result.Allowed {
				encoder(writer, request, srv.FormatError(request.Context(), result))
				return
			}
		}
		handler.ServeHTTP(writer, request)
	})
}
