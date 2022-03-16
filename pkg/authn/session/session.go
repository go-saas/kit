package session

import (
	"context"
	"github.com/goxiaoy/go-saas-kit/pkg/authn"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	"github.com/goxiaoy/sessions"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"net/http"
)

const (
	defaultSessionName  = "kit_user"
	defaultRememberName = "kit_user_rm"
)

func Auth(cfg *conf.Security) func(http.Handler) http.Handler {

	sessionInfoStore := NewSessionInfoStore(cfg)
	rememberStore := NewRememberStore(cfg)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := sessions.NewRegistryContext(r.Context(), r.Header)

			s, _ := GetSession(ctx, r.Header, sessionInfoStore, cfg)

			rs, _ := GetRememberSession(ctx, r.Header, rememberStore, cfg)

			stateWriter := NewClientStateWriter(s, rs, w, r.Header)

			ctx = NewClientStateWriterContext(ctx, stateWriter)
			state := NewClientState(s, rs)
			ctx = NewClientStateContext(ctx, state)
			ctx = authn.NewUserContext(ctx, authn.NewUserInfo(state.GetUid()))

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetSession(ctx context.Context, header sessions.Header, sessionInfoStore sessions.Store, cfg *conf.Security) (*sessions.Session, error) {
	var sn = defaultSessionName
	if cfg.SessionCookie != nil && cfg.SessionCookie.Name != nil {
		sn = cfg.SessionCookie.Name.Value
	}
	return sessionInfoStore.Get(ctx, header, sn)
}
func GetRememberSession(ctx context.Context, header sessions.Header, rememberStore sessions.Store, cfg *conf.Security) (*sessions.Session, error) {
	var rn = defaultRememberName
	if cfg.RememberCookie != nil && cfg.RememberCookie.Name != nil {
		rn = cfg.RememberCookie.Name.Value
	}
	return rememberStore.Get(ctx, header, rn)
}

//TODO handle remember?

func NewSessionInfoStore(cfg *conf.Security) sessions.Store {
	var blockKey []byte = nil
	if cfg.SecurityCookie.BlockKey != nil {
		blockKey = []byte(cfg.SecurityCookie.BlockKey.Value)
	}
	var store = sessions.NewCookieStore([]byte(cfg.SecurityCookie.HashKey), blockKey)
	if cfg.SessionCookie != nil {
		patchCfg(store, cfg.SessionCookie)
	}
	return store
}

func NewRememberStore(cfg *conf.Security) sessions.Store {
	var blockKey []byte = nil
	if cfg.SecurityCookie.BlockKey != nil {
		blockKey = []byte(cfg.SecurityCookie.BlockKey.Value)
	}
	var store = sessions.NewCookieStore([]byte(cfg.SecurityCookie.HashKey), blockKey)
	if cfg.RememberCookie == nil {
		cfg.RememberCookie = &conf.Cookie{}
	}
	if cfg.RememberCookie.MaxAge == nil {
		//365 days
		cfg.RememberCookie.MaxAge = &wrapperspb.Int32Value{Value: 86400 * 30 * 365}
	}
	patchCfg(store, cfg.RememberCookie)
	return store
}

func patchCfg(store *sessions.CookieStore, c *conf.Cookie) {
	if c.MaxAge != nil {
		store.MaxAge(int(c.MaxAge.Value))
	}
	if c.Path != nil {
		store.Options.Path = c.Path.Value
	}
	if c.HttpOnly != nil {
		store.Options.HttpOnly = c.HttpOnly.Value
	}
	if c.Secure != nil {
		store.Options.Secure = c.Secure.Value
	}
	if c.SameSite != conf.SameSiteMode_SameSiteNone {
		switch c.SameSite {
		case conf.SameSiteMode_SameSiteLax:
			store.Options.SameSite = http.SameSiteLaxMode
		case conf.SameSiteMode_SameSiteNone:
			store.Options.SameSite = http.SameSiteNoneMode
		case conf.SameSiteMode_SameSiteStrict:
			store.Options.SameSite = http.SameSiteStrictMode
		default:
			store.Options.SameSite = http.SameSiteDefaultMode
		}
	}
}
