package session

import (
	"context"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/goxiaoy/go-saas-kit/pkg/authn"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	"github.com/goxiaoy/go-saas-kit/pkg/errors"
	"github.com/goxiaoy/sessions"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"net/http"
)

const (
	defaultSessionName  = "kit_user"
	defaultRememberName = "kit_user_rm"
)

type RefreshTokenProvider = func(ctx context.Context, token string) (err error)

func Auth(cfg *conf.Security) func(http.Handler) http.Handler {

	sessionInfoStore := NewSessionInfoStore(cfg)
	rememberStore := NewRememberStore(cfg)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := sessions.NewRegistryContext(r.Context(), r.Header)

			s, _ := GetSession(ctx, r.Header, sessionInfoStore, cfg)

			rs, _ := GetRememberSession(ctx, r.Header, rememberStore, cfg)

			stateWriter := NewClientStateWriter(s, rs, w.Header(), r.Header)

			ctx = NewClientStateWriterContext(ctx, stateWriter)
			state := NewClientState(s, rs)
			ctx = NewClientStateContext(ctx, state)
			ctx = authn.NewUserContext(ctx, authn.NewUserInfo(state.GetUid()))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Refresh(errEncoder khttp.EncodeErrorFunc, provider RefreshTokenProvider) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			if state, ok := FromClientStateContext(ctx); ok {
				if stateWriter, ok := FromClientStateWriterContext(ctx); ok {
					if len(state.GetUid()) == 0 && len(state.GetRememberToken()) > 0 {
						//call refresh
						err := provider(ctx, state.GetRememberToken())
						if err != nil {
							if errors.Recoverable(err) {
								//abort with error
								errEncoder(w, r, err)
								return
							} else {
								//just clean remember token
								stateWriter.SignOutRememberToken(ctx)
								stateWriter.Save(ctx)
							}
						} else {
							ctx = authn.NewUserContext(ctx, authn.NewUserInfo(state.GetUid()))
						}
					}
				}
			}
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

func NewSessionInfoStore(cfg *conf.Security) sessions.Store {
	var blockKey []byte = nil
	if cfg.SecurityCookie.BlockKey != nil {
		blockKey = []byte(cfg.SecurityCookie.BlockKey.Value)
	}
	var store = sessions.NewCookieStore([]byte(cfg.SecurityCookie.HashKey), blockKey)

	if cfg.SessionCookie == nil {
		cfg.SessionCookie = &conf.Cookie{}
	}
	if cfg.SessionCookie.MaxAge == nil {
		cfg.SessionCookie.MaxAge = &wrapperspb.Int32Value{Value: int32(SessionExpireSecondsOrDefault(0))}
	}
	patchCfg(store, cfg.SessionCookie)

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
		cfg.RememberCookie.MaxAge = &wrapperspb.Int32Value{Value: int32(RememberMeExpireSecondsOrDefault(0))}
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

//SessionExpireSecondsOrDefault default 1 day
func SessionExpireSecondsOrDefault(seconds int) int {
	if seconds <= 0 {
		//1 day
		return 24 * 60 * 60
	}
	return seconds
}

//RememberMeExpireSecondsOrDefault default 30 days
func RememberMeExpireSecondsOrDefault(seconds int) int {
	if seconds <= 0 {
		//30 days
		return 24 * 60 * 60 * 30
	}
	return seconds
}
