package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gorilla/sessions"
	conf2 "github.com/goxiaoy/go-saas-kit/pkg/conf"
	sms2 "github.com/goxiaoy/go-saas-kit/pkg/sms"
	"github.com/goxiaoy/go-saas-kit/user/private/biz"
	"github.com/goxiaoy/go-saas-kit/user/private/conf"
	abclientstate "github.com/volatiletech/authboss-clientstate"
	"github.com/volatiletech/authboss-renderer"
	"github.com/volatiletech/authboss/v3"
	"github.com/volatiletech/authboss/v3/defaults"
	_ "github.com/volatiletech/authboss/v3/logout"
	"github.com/volatiletech/authboss/v3/otp/twofactor"
	"github.com/volatiletech/authboss/v3/otp/twofactor/sms2fa"
	"github.com/volatiletech/authboss/v3/otp/twofactor/totp2fa"
	_ "github.com/volatiletech/authboss/v3/recover"
	_ "github.com/volatiletech/authboss/v3/register"
	"net/http"
	"regexp"
)

const (
	sessionCookieName = "kit_user"
	totp2FAIssuer     = "kit_user"
	apiMode           = true
)

func NewSessionStorer(sCfg *conf2.Security, u *conf.UserConf) *abclientstate.SessionStorer {
	cfg := u.Auth
	sn := sessionCookieName
	if cfg.SessionName != nil {
		sn = cfg.SessionName.Value
	}
	var blockKey []byte = nil
	if sCfg.SecurityCookie.BlockKey != nil {
		blockKey = []byte(sCfg.SecurityCookie.BlockKey.Value)
	}
	sessionStore := abclientstate.NewSessionStorer(sn, []byte(sCfg.SecurityCookie.HashKey), blockKey)
	cstore := sessionStore.Store.(*sessions.CookieStore)
	if cfg.SessionCookie != nil {
		c := cfg.SessionCookie
		if c.MaxAge != nil {
			cstore.MaxAge(int(c.MaxAge.Value))
		}
		if c.Path != nil {
			cstore.Options.Path = c.Path.Value
		}
		if c.HttpOnly != nil {
			cstore.Options.HttpOnly = c.HttpOnly.Value
		}
		if c.Secure != nil {
			cstore.Options.Secure = c.Secure.Value
		}
		if c.SameSite != conf2.SameSiteMode_SameSiteNone {
			switch c.SameSite {
			case conf2.SameSiteMode_SameSiteLax:
				cstore.Options.SameSite = http.SameSiteLaxMode
			case conf2.SameSiteMode_SameSiteNone:
				cstore.Options.SameSite = http.SameSiteNoneMode
			case conf2.SameSiteMode_SameSiteStrict:
				cstore.Options.SameSite = http.SameSiteStrictMode
			default:
				cstore.Options.SameSite = http.SameSiteDefaultMode
			}
		}
	}
	return &sessionStore
}

func NewCookieStorer(sCfg *conf2.Security, u *conf.UserConf) *abclientstate.CookieStorer {
	cfg := u.Auth
	var blockKey []byte = nil
	if sCfg.SecurityCookie.BlockKey != nil {
		blockKey = []byte(sCfg.SecurityCookie.BlockKey.Value)
	}
	cookieStore := abclientstate.NewCookieStorer([]byte(sCfg.SecurityCookie.HashKey), blockKey)
	if cfg.Cookie != nil {
		c := cfg.Cookie
		if c.MaxAge != nil {
			cookieStore.MaxAge = int(c.MaxAge.Value)
		}
		if c.Domain != nil {
			cookieStore.Domain = c.Domain.Value
		}
		if c.Path != nil {
			cookieStore.Path = c.Path.Value
		}
		if c.HttpOnly != nil {
			cookieStore.HTTPOnly = c.HttpOnly.Value
		}
		if c.Secure != nil {
			cookieStore.Secure = c.Secure.Value
		}
		if c.SameSite != conf2.SameSiteMode_SameSiteNone {
			switch c.SameSite {
			case conf2.SameSiteMode_SameSiteLax:
				cookieStore.SameSite = http.SameSiteLaxMode
			case conf2.SameSiteMode_SameSiteNone:
				cookieStore.SameSite = http.SameSiteNoneMode
			case conf2.SameSiteMode_SameSiteStrict:
				cookieStore.SameSite = http.SameSiteStrictMode
			default:
				cookieStore.SameSite = http.SameSiteDefaultMode
			}
		}
	}
	return &cookieStore
}

type logWrapper struct {
	*log.Helper
}

func (l *logWrapper) Info(s string) {
	l.Helper.Info(s)
}
func (l *logWrapper) Error(s string) {
	l.Helper.Error(s)
}

func NewAuthboss(l log.Logger, u *conf.UserConf, session *abclientstate.SessionStorer, cookie *abclientstate.CookieStorer, store *biz.AuthbossStoreWrapper) (*authboss.Authboss, error) {
	ab := authboss.New()

	ab.Config.Storage.Server = store
	ab.Config.Storage.SessionState = session
	ab.Config.Storage.CookieState = cookie
	ab.Config.Paths.Mount = "/v1/auth/web"
	ab.Config.Paths.RootURL = u.RootUrl

	logger := log.NewHelper(l)

	ab.Config.Core.Logger = &logWrapper{
		logger,
	}
	// This is using the renderer from: github.com/volatiletech/authboss
	ab.Config.Core.ViewRenderer = defaults.JSONRenderer{}
	// We render mail with the authboss-renderer but we use a LogMailer
	// which simply sends the e-mail to stdout.
	ab.Config.Core.MailRenderer = abrenderer.NewEmail("/auth", "ab_views")

	// The preserve fields are things we don't want to
	// lose when we're doing user registration (prevents having
	// to type them again)
	ab.Config.Modules.RegisterPreserveFields = []string{"email", "name"}

	// TOTP2FAIssuer is the name of the issuer we use for totp 2fa
	TOTP2FAIssuer := totp2FAIssuer
	if u.Auth.Totp_2FaIssuer != nil {
		TOTP2FAIssuer = u.Auth.Totp_2FaIssuer.Value
	}
	ab.Config.Modules.TOTP2FAIssuer = TOTP2FAIssuer
	ab.Config.Modules.ResponseOnUnauthed = authboss.RespondRedirect

	// Turn on e-mail authentication required
	ab.Config.Modules.TwoFactorEmailAuthRequired = true

	// This instantiates and uses every default implementation
	// in the Config.Core area that exist in the defaults package.
	// Just a convenient helper if you don't want to do anything fancy.
	defaults.SetCore(&ab.Config, apiMode, false)

	// Here we initialize the bodyreader as something customized in order to accept a name
	// parameter for our user as well as the standard e-mail and password.
	//
	// We also change the validation for these fields
	// to be something less secure so that we can use test data easier.
	emailRule := defaults.Rules{
		FieldName: "email", Required: true,
		MatchError: "Must be a valid e-mail address",
		MustMatch:  regexp.MustCompile(`.*@.*\.[a-z]+`),
	}
	passwordRule := defaults.Rules{
		FieldName: "password", Required: true,
		MinLength: 4,
	}
	nameRule := defaults.Rules{
		FieldName: "name", Required: true,
		MinLength: 2,
	}

	ab.Config.Core.BodyReader = defaults.HTTPBodyReader{
		ReadJSON: apiMode,
		Rulesets: map[string][]defaults.Rules{
			"register":    {emailRule, passwordRule, nameRule},
			"recover_end": {passwordRule},
		},
		Confirms: map[string][]string{
			"register":    {"password", authboss.ConfirmPrefix + "password"},
			"recover_end": {"password", authboss.ConfirmPrefix + "password"},
		},
		Whitelist: map[string][]string{
			"register": {"email", "name", "password"},
		},
	}

	// Set up 2fa
	twofaRecovery := &twofactor.Recovery{Authboss: ab}
	if err := twofaRecovery.Setup(); err != nil {
		panic(err)
	}

	totp := &totp2fa.TOTP{Authboss: ab}
	if err := totp.Setup(); err != nil {
		panic(err)
	}

	sms := &sms2fa.SMS{Authboss: ab, Sender: &sms2.SMSLogSender{}}
	if err := sms.Setup(); err != nil {
		panic(err)
	}

	if err := ab.Init(); err != nil {
		return nil, err
	}
	return ab, nil
}
