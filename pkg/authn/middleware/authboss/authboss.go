package authboss

import (
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/goxiaoy/go-saas-kit/pkg/authn"
	"github.com/volatiletech/authboss/v3"
	"github.com/volatiletech/authboss/v3/remember"
	"net/http"
)

func PathFilter(ab *authboss.Authboss) func(http.Handler) http.Handler {
	//1. load client state
	//2. hook system user
	//3. handle remember me
	chain := khttp.FilterChain(ab.LoadClientStateMiddleware, hookState(ab), rememberMiddleware(ab))
	return chain
}

func hookState(ab *authboss.Authboss) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if id, _ := ab.CurrentUserID(r); len(id) != 0 {
				r.WithContext(authn.NewUserContext(r.Context(), authn.NewUserInfo(id)))
			}
			next.ServeHTTP(w, r)
		})
	}
}

// rememberMiddleware automatically authenticates users if they have remember me tokens
// If the user has been loaded already, it returns early
// see remember.Middleware
func rememberMiddleware(ab *authboss.Authboss) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Safely can ignore error here
			if user, ok := authn.FromUserContext(r.Context()); !ok || len(user.GetId()) == 0 {
				//refresh remember
				if err := remember.Authenticate(ab, w, &r); err != nil {
					logger := ab.RequestLogger(r)
					logger.Errorf("failed to authenticate user via remember me: %+v", err)
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
