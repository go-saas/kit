package csrf

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	http2 "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-saas/kit/pkg/conf"
	"github.com/gorilla/csrf"
	"net/http"
)

func NewCsrf(l log.Logger, sCfg *conf.Security, cfg *conf.Server_HTTP_Csrf, errEncoder http2.EncodeErrorFunc) func(http.Handler) http.Handler {
	logger := log.NewHelper(l)

	var csrfOpt []csrf.Option

	name := "kit_csrf"

	if cfg.Cookie != nil {
		if cfg.Cookie.Name != nil {
			name = cfg.Cookie.Name.Value
		}
		if cfg.Cookie.MaxAge != nil {
			csrfOpt = append(csrfOpt, csrf.MaxAge(int(cfg.Cookie.MaxAge.Value)))
		}
		if cfg.Cookie.Domain != nil {
			csrfOpt = append(csrfOpt, csrf.Domain(cfg.Cookie.Domain.Value))
		}
		csrfOpt = append(csrfOpt, csrf.Path("/"))
		if cfg.Cookie.HttpOnly != nil {
			csrfOpt = append(csrfOpt, csrf.HttpOnly(cfg.Cookie.HttpOnly.Value))
		}
		if cfg.Cookie.Secure != nil {
			csrfOpt = append(csrfOpt, csrf.Secure(cfg.Cookie.Secure.Value))
		}
		if cfg.Cookie.SameSite != conf.SameSiteMode_SameSiteDefault {
			switch cfg.Cookie.SameSite {
			case conf.SameSiteMode_SameSiteLax:
				csrfOpt = append(csrfOpt, csrf.SameSite(csrf.SameSiteLaxMode))
			case conf.SameSiteMode_SameSiteNone:
				csrfOpt = append(csrfOpt, csrf.SameSite(csrf.SameSiteNoneMode))
			case conf.SameSiteMode_SameSiteStrict:
				csrfOpt = append(csrfOpt, csrf.SameSite(csrf.SameSiteStrictMode))
			}
		}
	}
	csrfOpt = append(csrfOpt, csrf.CookieName(name))
	if cfg.RequestHeader != nil {
		csrfOpt = append(csrfOpt, csrf.RequestHeader(cfg.RequestHeader.Value))
	}
	if cfg.FieldName != nil {
		csrfOpt = append(csrfOpt, csrf.FieldName(cfg.FieldName.Value))
	}

	if len(cfg.TrustedOrigins) > 0 {
		csrfOpt = append(csrfOpt, csrf.TrustedOrigins(cfg.TrustedOrigins))
	}

	// unauthorizedhandler sets a HTTP 403 Forbidden status and writes the
	// CSRF failure reason to the response.
	unauthorizedHandler := func(w http.ResponseWriter, r *http.Request) {
		logger.Debugf(fmt.Sprintf("%s - %s",
			http.StatusText(http.StatusForbidden), csrf.FailureReason(r)))

		errEncoder(w, r, errors.New(http.StatusForbidden, "CSRF_INVALID", csrf.FailureReason(r).Error()))

		return
	}
	csrfOpt = append(csrfOpt, csrf.ErrorHandler(http.HandlerFunc(unauthorizedHandler)))
	CSRF := csrf.Protect([]byte(sCfg.SecurityCookie.HashKey), csrfOpt...)

	return CSRF
}
