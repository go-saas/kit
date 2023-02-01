package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-saas/kit/pkg/authn"
	"github.com/go-saas/kit/pkg/authn/jwt"
	"github.com/go-saas/kit/pkg/authz/authz"
	"github.com/go-saas/kit/pkg/conf"
	"github.com/go-saas/kit/pkg/csrf"
	"github.com/go-saas/kit/pkg/server/endpoint"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/goxiaoy/vfs"
	"github.com/spf13/afero"
	"net"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type (
	// ServiceRegister register http handler into http server
	ServiceRegister interface {
		Register(server *khttp.Server, middleware ...middleware.Middleware)
	}
	ServiceRegisterFunc func(server *khttp.Server, middleware ...middleware.Middleware)
)

var (
	ReqDecode  khttp.DecodeRequestFunc  = khttp.DefaultRequestDecoder
	ResEncoder khttp.EncodeResponseFunc = khttp.DefaultResponseEncoder
	ErrEncoder khttp.EncodeErrorFunc    = khttp.DefaultErrorEncoder
)

func (f ServiceRegisterFunc) Register(server *khttp.Server, middleware ...middleware.Middleware) {
	f(server, middleware...)
}

func ChainServiceRegister(r ...ServiceRegister) ServiceRegister {
	return ServiceRegisterFunc(func(server *khttp.Server, middleware ...middleware.Middleware) {
		for _, register := range r {
			register.Register(server, middleware...)
		}
	})
}

// PatchOpts Patch http options with given service name and configs. f use global filters
func PatchOpts(l log.Logger,
	opts []khttp.ServerOption,
	server *conf.Server,
	sCfg *conf.Security,
	reqDecoder khttp.DecodeRequestFunc,
	resEncoder khttp.EncodeResponseFunc,
	errEncoder khttp.EncodeErrorFunc,
	f ...khttp.FilterFunc) []khttp.ServerOption {

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
			handlers.AllowedHeaders(append([]string{"Content-Type", jwt.AuthorizationHeader}, server.Http.Cors.AllowedHeaders...)),
		))
	}
	if server.Http.Csrf != nil {
		filters = append(filters, csrf.NewCsrf(l, sCfg, server.Http.Csrf, errEncoder))
	}
	filters = append(filters, f...)
	opts = append(opts, khttp.Filter(filters...))
	return opts
}

func MountBlob(srv *khttp.Server, pathPrefix, basePath string, b vfs.Blob) {
	if pathPrefix == "" {
		pathPrefix = "/assets"
	}
	if !strings.HasPrefix(pathPrefix, "/") {
		pathPrefix = fmt.Sprintf("/%s", pathPrefix)
	}
	fullPath := path.Join(pathPrefix, basePath)
	router := mux.NewRouter()
	//expose

	router.PathPrefix(fullPath).Handler(http.StripPrefix(fullPath, http.FileServer(http.FS(afero.NewIOFS(afero.NewBasePathFs(b, basePath))))))

	srv.HandlePrefix(fullPath, router)
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

func SetCookie(ctx context.Context, cookie *http.Cookie) error {
	if t, ok := transport.FromServerContext(ctx); ok {
		if v := cookie.String(); v != "" {
			t.ReplyHeader().Set("Set-Cookie", v)
		}
		return nil
	} else {
		return errors.New("unsupported transport")
	}
}

func AuthnGuardian(encoder khttp.EncodeErrorFunc, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, err := authn.ErrIfUnauthenticated(request.Context())
		if err != nil {
			encoder(writer, request, err)
			return
		}
		handler.ServeHTTP(writer, request)
	})
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
				encoder(writer, request, srv.FormatError(request.Context(), requirement, result))
				return
			}
		}
		handler.ServeHTTP(writer, request)
	})
}

type Server struct {
	*khttp.Server
	cfg *conf.Server
}

func NewServer(cfg *conf.Server, opts ...khttp.ServerOption) *Server {
	return &Server{
		Server: khttp.NewServer(opts...),
		cfg:    cfg,
	}
}

func (s *Server) Endpoint() (url *url.URL, err error) {
	url, err = s.Server.Endpoint()
	if err != nil || url == nil {
		return
	}
	if s.cfg.Http != nil && len(s.cfg.Http.Endpoint) > 0 {
		//TODO tls
		return endpoint.NewEndpoint(endpoint.Scheme("http", false), s.cfg.Http.Endpoint), nil
	}
	return
}
