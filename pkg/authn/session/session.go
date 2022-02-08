package session

import (
	"context"
	"github.com/gorilla/sessions"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	"net/http"
)

const (
	defaultSessionName = "kit-user"
)

func Auth(store sessions.Store, cfg *conf.Security) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var sn = defaultSessionName
			if cfg.GetSessionName() != nil {
				sn = cfg.GetSessionName().GetValue()
			}
			s, _ := store.Get(r, sn)
			stateWriter := NewClientStateWriter(s, w, r)
			defer func() { stateWriter.Save(context.Background()) }()
			newCtx := NewClientStateWriterContext(r.Context(), stateWriter)

			if s != nil {
				state := NewClientState(s)

				newCtx = NewClientStateContext(newCtx, state)
				next.ServeHTTP(w, r.WithContext(newCtx))
			} else {
				next.ServeHTTP(w, r.WithContext(newCtx))
			}

		})
	}

}

func NewCookieStore(cfg *conf.Security) sessions.Store {
	var blockKey []byte = nil
	if cfg.SecurityCookie.BlockKey != nil {
		blockKey = []byte(cfg.SecurityCookie.BlockKey.Value)
	}
	var store = sessions.NewCookieStore([]byte(cfg.SecurityCookie.HashKey), blockKey)
	if cfg.SessionCookie != nil {
		c := cfg.SessionCookie
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
	return store
}
