package http

import (
	"github.com/go-kratos/kratos/v2/log"
	http2 "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/kratos/v2/transport/http/binding"
	v1 "github.com/goxiaoy/go-saas-kit/user/api/auth/v1"
	"github.com/goxiaoy/go-saas-kit/user/private/biz"
	"github.com/goxiaoy/go-saas-kit/user/private/service"
	"github.com/volatiletech/authboss/v3"
	"net/http"
)

type Auth struct {
	*authboss.Authboss
	reqDecoder http2.DecodeRequestFunc
	resEncoder http2.EncodeResponseFunc
	um         *biz.UserManager
	logger     *log.Helper
}

func NewAuth(authboss *authboss.Authboss,
	reqDecoder http2.DecodeRequestFunc,
	resEncoder http2.EncodeResponseFunc,
	um *biz.UserManager,
	l log.Logger) *Auth {
	return &Auth{
		authboss,
		reqDecoder,
		resEncoder,
		um,
		log.NewHelper(l),
	}
}

func (a *Auth) LoginGet(w http.ResponseWriter, r *http.Request) error {
	var req v1.GetLoginFormRequest
	if err := binding.BindQuery(r.URL.Query(), &req); err != nil {
		return err
	}
	var resp v1.GetLoginFormResponse
	//TODO validate url
	resp.Redirect = req.Redirect
	resp.Oauth = make([]*v1.OAuthProvider, len(a.Config.Modules.OAuth2Providers))
	for k, _ := range a.Config.Modules.OAuth2Providers {
		resp.Oauth = append(resp.Oauth, &v1.OAuthProvider{Name: k})
	}
	if err := a.resEncoder(w, r, &resp); err != nil {
		return err
	}
	return nil
}

func (a *Auth) LoginPost(w http.ResponseWriter, r *http.Request) error {
	//find user
	var req v1.LoginAuthRequest
	if err := a.reqDecoder(r, &req); err != nil {
		return err
	}
	var handled bool
	user, err := service.FindUserByUsernameAndValidatePwd(r.Context(), a.um, req.Username, req.Password)

	if err != nil {
		handled, err = a.Authboss.Events.FireAfter(authboss.EventAuthFail, w, r)
		if err != nil {
			return err
		} else if handled {
			return nil
		}
		a.logger.Infof("user with username %s failed to log in", req.Username)
		return err
	}

	//TODO
	//r = r.WithContext(context.WithValue(r.Context(), authboss.CTXKeyValues, validatable))

	handled, err = a.Events.FireBefore(authboss.EventAuth, w, r)
	if err != nil {
		return err
	} else if handled {
		return nil
	}

	handled, err = a.Events.FireBefore(authboss.EventAuthHijack, w, r)
	if err != nil {
		return err
	} else if handled {
		return nil
	}

	a.logger.Infof("user %s logged in", user.ID.String())
	authboss.PutSession(w, authboss.SessionKey, user.ID.String())
	authboss.DelSession(w, authboss.SessionHalfAuthKey)

	handled, err = a.Authboss.Events.FireAfter(authboss.EventAuth, w, r)
	if err != nil {
		return err
	} else if handled {
		return nil
	}

	ro := authboss.RedirectOptions{
		Code:             http.StatusTemporaryRedirect,
		RedirectPath:     a.Authboss.Paths.AuthLoginOK,
		FollowRedirParam: true,
	}
	return a.Authboss.Core.Redirector.Redirect(w, r, ro)
}
