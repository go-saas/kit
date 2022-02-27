package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/goxiaoy/go-saas-kit/pkg/blob"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	"github.com/goxiaoy/go-saas-kit/pkg/csrf"
	"github.com/spf13/afero"
	"google.golang.org/protobuf/proto"
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

const defaultSrvName = "default"

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
	var server *conf.Server
	if s, ok := services.Servers[name]; ok {
		server = s
	}
	if def, ok := services.Servers[defaultSrvName]; ok {
		if server != nil {
			proto.Merge(server, def)
		} else {
			server = def
		}
	} else if server == nil {
		panic(errors.New(fmt.Sprintf("both %v and %s server not found", name, defaultSrvName)))
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
