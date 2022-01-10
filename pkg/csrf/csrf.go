package csrf

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	http2 "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/csrf"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	"net/http"
)

func NewCsrf(l log.Logger, cfg *conf.Server_HTTP_Csrf) func(http.Handler) http.Handler {
	logger := log.NewHelper(l)

	var csrfOpt []csrf.Option

	if cfg.Cookie != nil {
		if cfg.Cookie.Name != nil {
			csrfOpt = append(csrfOpt, csrf.CookieName(cfg.Cookie.Name.Value))
		}
		if cfg.Cookie.MaxAge != nil {
			csrfOpt = append(csrfOpt, csrf.MaxAge(int(cfg.Cookie.MaxAge.Value)))
		}
		if cfg.Cookie.Domain != nil {
			csrfOpt = append(csrfOpt, csrf.Domain(cfg.Cookie.Domain.Value))
		}
		if cfg.Cookie.Path != nil {
			csrfOpt = append(csrfOpt, csrf.Path(cfg.Cookie.Path.Value))
		}
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

		//TODO possible replacement of error encoder
		http2.DefaultErrorEncoder(w, r, errors.New(http.StatusForbidden, "CSRF_INVALID", csrf.FailureReason(r).Error()))

		return
	}
	csrfOpt = append(csrfOpt, csrf.ErrorHandler(http.HandlerFunc(unauthorizedHandler)))
	CSRF := csrf.Protect([]byte(cfg.AuthKey), csrfOpt...)

	return CSRF
}
