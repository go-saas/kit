package http

import (
	"github.com/go-kratos/kratos/v2/log"
	http2 "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/kratos/v2/transport/http/binding"
	v1 "github.com/goxiaoy/go-saas-kit/user/api/auth/v1"
	"github.com/goxiaoy/go-saas-kit/user/private/biz"
	"github.com/goxiaoy/go-saas-kit/user/private/service"
	"net/http"
)

type Auth struct {
	reqDecoder http2.DecodeRequestFunc
	resEncoder http2.EncodeResponseFunc
	um         *biz.UserManager
	logger     *log.Helper
	signIn     *biz.SignInManager
}

func NewAuth(
	reqDecoder http2.DecodeRequestFunc,
	resEncoder http2.EncodeResponseFunc,
	um *biz.UserManager,
	l log.Logger,
	signIn *biz.SignInManager) *Auth {
	return &Auth{
		reqDecoder,
		resEncoder,
		um,
		log.NewHelper(l),
		signIn,
	}
}

func (a *Auth) LoginGet(w http.ResponseWriter, r *http.Request) error {
	var req v1.GetLoginRequest
	if err := binding.BindQuery(r.URL.Query(), &req); err != nil {
		return err
	}

	var resp v1.GetLoginResponse
	//TODO validate url
	resp.Redirect = req.Redirect
	//TODO oauth
	//resp.Oauth = make([]*v1.OAuthProvider, len(a.Config.Modules.OAuth2Providers))
	//for k, _ := range a.Config.Modules.OAuth2Providers {
	//	resp.Oauth = append(resp.Oauth, &v1.OAuthProvider{Name: k})
	//}
	if err := a.resEncoder(w, r, &resp); err != nil {
		return err
	}
	return nil
}

func (a *Auth) LoginPost(w http.ResponseWriter, r *http.Request) error {
	var req v1.LoginAuthRequest
	if err := a.reqDecoder(r, &req); err != nil {
		return err
	}
	err := a.signIn.PasswordSignInWithUsername(r.Context(), req.Username, req.Password, req.Remember, true)
	return service.ConvertError(err)
}
func (a *Auth) LoginOut(w http.ResponseWriter, r *http.Request) error {
	var req v1.LogoutRequest
	if err := a.reqDecoder(r, &req); err != nil {
		return err
	}
	err := a.signIn.SignOut(r.Context())
	return service.ConvertError(err)
}
