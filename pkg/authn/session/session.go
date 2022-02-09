package session

import (
	"context"
	"github.com/gorilla/sessions"
	"github.com/goxiaoy/go-saas-kit/pkg/authn"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
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
			var sn = defaultSessionName
			if cfg.SessionCookie != nil && cfg.SessionCookie.Name != nil {
				sn = cfg.SessionCookie.Name.Value
			}

			var rn = defaultRememberName
			if cfg.RememberCookie != nil && cfg.RememberCookie.Name != nil {
				rn = cfg.RememberCookie.Name.Value
			}

			s, _ := sessionInfoStore.Get(r, sn)

			rs, _ := rememberStore.Get(r, rn)

			stateWriter := NewClientStateWriter(s, rs, w, r)
			defer func() { stateWriter.Save(context.Background()) }()
			newCtx := NewClientStateWriterContext(r.Context(), stateWriter)
			state := NewClientState(s, rs)
			newCtx = NewClientStateContext(newCtx, state)
			newCtx = authn.NewUserContext(newCtx, authn.NewUserInfo(state.GetUid()))

			next.ServeHTTP(w, r.WithContext(newCtx))

		})
	}
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
