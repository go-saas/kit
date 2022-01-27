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
}

func NewAuth(
	reqDecoder http2.DecodeRequestFunc,
	resEncoder http2.EncodeResponseFunc,
	um *biz.UserManager,
	l log.Logger) *Auth {
	return &Auth{
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
	//find user
	var req v1.LoginAuthRequest
	if err := a.reqDecoder(r, &req); err != nil {
		return err
	}
	var handled bool
	_, err := service.FindUserByUsernameAndValidatePwd(r.Context(), a.um, req.Username, req.Password)

	if err != nil {
		if err != nil {
			return err
		} else if handled {
			return nil
		}
		a.logger.Infof("user with username %s failed to log in", req.Username)
		return err
	}

	//TODO
	panic("not implementd")

}
