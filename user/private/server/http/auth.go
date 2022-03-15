package http

import (
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	http2 "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/kratos/v2/transport/http/binding"
	v1 "github.com/goxiaoy/go-saas-kit/user/api/auth/v1"
	"github.com/goxiaoy/go-saas-kit/user/private/biz"
	"github.com/goxiaoy/go-saas-kit/user/private/service"
	"github.com/ory/hydra-client-go"
	"net/http"
)

type Auth struct {
	reqDecoder http2.DecodeRequestFunc
	resEncoder http2.EncodeResponseFunc
	um         *biz.UserManager
	logger     *log.Helper
	signIn     *biz.SignInManager
	hclient    *client.APIClient
}

func NewAuth(
	reqDecoder http2.DecodeRequestFunc,
	resEncoder http2.EncodeResponseFunc,
	um *biz.UserManager,
	l log.Logger,
	signIn *biz.SignInManager,
	hclient *client.APIClient,
) *Auth {
	return &Auth{
		reqDecoder: reqDecoder,
		resEncoder: resEncoder,
		um:         um,
		logger:     log.NewHelper(l),
		signIn:     signIn,
		hclient:    hclient,
	}
}

func (a *Auth) LoginGet(w http.ResponseWriter, r *http.Request) error {
	var req v1.GetLoginRequest
	if err := binding.BindQuery(r.URL.Query(), &req); err != nil {
		return err
	}
	//TODO already signin?
	var resp v1.GetLoginResponse
	resp.Redirect = req.Redirect
	if len(req.LoginChallenge) > 0 {
		//hydra
		if hreq, _, err := a.hclient.AdminApi.GetLoginRequestExecute(a.hclient.AdminApi.GetLoginRequest(r.Context())); err != nil {
			return err
		} else {
			// If hydra was already able to authenticate the user, skip will be true and we do not need to re-authenticate
			// the user.
			if hreq.Skip {
				acc, _, err := a.hclient.AdminApi.AcceptLoginRequest(r.Context()).
					LoginChallenge(hreq.Challenge).
					AcceptLoginRequest(*client.NewAcceptLoginRequest(hreq.Subject)).
					Execute()
				if err != nil {
					return err
				}
				resp.Redirect = acc.RedirectTo
			} else {
				resp.Challenge = hreq.Challenge
				if hreq.OidcContext != nil && hreq.OidcContext.LoginHint != nil {
					resp.Hint = *hreq.OidcContext.LoginHint
				}
			}
		}
	}

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
	var req v1.WebLoginAuthRequest
	if err := a.reqDecoder(r, &req); err != nil {
		return err
	}
	var resp v1.WebLoginAuthReply
	if len(req.Challenge) > 0 {
		// Let's see if the user decided to accept or reject the consent request..
		if req.Reject {
			// Looks like the consent request was denied by the user
			reject := *client.NewRejectRequest()
			//TODO error
			reject.SetError("access_denied")
			reject.SetErrorDescription("The resource owner denied the request")
			hreq, _, err := a.hclient.AdminApi.RejectLoginRequest(r.Context()).LoginChallenge(req.Challenge).RejectRequest(reject).Execute()
			if err != nil {
				return err
			}
			resp.Redirect = hreq.RedirectTo
			//return
			if err := a.resEncoder(w, r, &resp); err != nil {
				return err
			}
			return nil
		}
	}
	//validate sign in
	err, uid := a.signIn.PasswordSignInWithUsername(r.Context(), req.Username, req.Password, req.Remember, true)
	if err != nil {
		return service.ConvertError(err)
	}
	if len(req.Challenge) > 0 {
		// Seems like the user authenticated! Let's tell hydra...
		acc := *client.NewAcceptLoginRequest(uid)
		acc.SetRemember(req.Remember)
		//TODO from config
		acc.SetRememberFor(3600)
		hreq, _, err := a.hclient.AdminApi.AcceptLoginRequest(r.Context()).
			LoginChallenge(req.Challenge).
			AcceptLoginRequest(acc).Execute()
		if err != nil {
			return err
		}
		resp.Redirect = hreq.RedirectTo
	}
	if err := a.resEncoder(w, r, &resp); err != nil {
		return err
	}
	return nil
}

func (a *Auth) LoginOutGet(w http.ResponseWriter, r *http.Request) error {
	var req v1.GetLogoutRequest
	if err := binding.BindQuery(r.URL.Query(), &req); err != nil {
		return err
	}
	var resp v1.GetLogoutResponse
	if len(req.LogoutChallenge) > 0 {
		_, _, err := a.hclient.AdminApi.GetLogoutRequest(r.Context()).LogoutChallenge(req.LogoutChallenge).Execute()
		if err != nil {
			return err
		}
		resp.Challenge = req.LogoutChallenge
	}
	if err := a.resEncoder(w, r, &resp); err != nil {
		return err
	}
	return nil
}

func (a *Auth) Logout(w http.ResponseWriter, r *http.Request) error {
	var req v1.LogoutRequest
	if err := a.reqDecoder(r, &req); err != nil {
		return err
	}
	var resp v1.LogoutResponse
	err := a.signIn.SignOut(r.Context())
	if err != nil {
		return service.ConvertError(err)
	}
	if len(req.Challenge) > 0 {
		hreq, _, err := a.hclient.AdminApi.AcceptLogoutRequest(r.Context()).LogoutChallenge(req.Challenge).Execute()
		if err != nil {
			return err
		}
		resp.Redirect = hreq.RedirectTo
	}
	if err := a.resEncoder(w, r, &resp); err != nil {
		return err
	}
	return nil

}

func (a *Auth) ConsentGet(w http.ResponseWriter, r *http.Request) error {
	var req v1.GetConsentRequest
	if err := binding.BindQuery(r.URL.Query(), &req); err != nil {
		return err
	}
	var resp v1.GetConsentResponse
	if len(req.ConsentChallenge) == 0 {
		return errors.BadRequest("CONSENT_CHALLENGE_REQUIRED", "")
	}
	hreq, _, err := a.hclient.AdminApi.GetConsentRequest(r.Context()).ConsentChallenge(req.ConsentChallenge).Execute()
	if err != nil {
		return err
	}
	if hreq.GetSkip() {
		acc := *client.NewAcceptConsentRequest()
		acc.SetGrantScope(hreq.RequestedScope)
		acc.SetGrantAccessTokenAudience(hreq.RequestedAccessTokenAudience)
		//acc.SetSession(client.ConsentRequestSession{
		//	AccessToken: nil,
		//	IdToken:     nil,
		//})
		accReq, _, err := a.hclient.AdminApi.AcceptConsentRequest(r.Context()).ConsentChallenge(req.ConsentChallenge).AcceptConsentRequest(acc).Execute()
		if err != nil {
			return err
		}
		resp.Redirect = accReq.RedirectTo
		if err := a.resEncoder(w, r, &resp); err != nil {
			return err
		}
		return nil
	}
	resp.Challenge = hreq.Challenge
	resp.RequestedScope = hreq.RequestedScope
	resp.UserId = hreq.GetSubject()
	c := hreq.GetClient()
	resp.ClientId = c.GetClientId()
	if err := a.resEncoder(w, r, &resp); err != nil {
		return err
	}
	return nil
}
func (a *Auth) Consent(w http.ResponseWriter, r *http.Request) error {
	var req v1.GrantConsentRequest
	if err := a.reqDecoder(r, &req); err != nil {
		return err
	}
	if len(req.Challenge) == 0 {
		return errors.BadRequest("CHALLENGE_REQUIRED", "")
	}

	var resp v1.GrantConsentResponse

	if req.Reject {
		reject := *client.NewRejectRequest()
		//TODO
		reject.SetError("access_denied")
		reject.SetErrorDescription("The resource owner denied the request")
		hreq, _, err := a.hclient.AdminApi.RejectConsentRequest(r.Context()).ConsentChallenge(req.Challenge).RejectRequest(reject).Execute()
		if err != nil {
			return err
		}
		resp.Redirect = hreq.RedirectTo
		//return
		if err := a.resEncoder(w, r, &resp); err != nil {
			return err
		}
		return nil
	}
	//user allow
	// The session allows us to set session data for id and access tokens
	session := client.NewConsentRequestSession()
	// This data will be available when introspecting the token. Try to avoid sensitive information here,
	// unless you limit who can introspect tokens.
	session.SetAccessToken(map[string]map[string]interface{}{})

	session.SetIdToken(map[string]map[string]interface{}{})

	// Here is also the place to add data to the ID or access token. For example,
	// if the scope 'profile' is added, add the family and given name to the ID Token claims:
	// if (grantScope.indexOf('profile')) {
	//   session.id_token.family_name = 'Doe'
	//   session.id_token.given_name = 'John'
	// }

	hreq, _, err := a.hclient.AdminApi.GetConsentRequest(r.Context()).ConsentChallenge(req.Challenge).Execute()
	if err != nil {
		return err
	}

	acc := client.NewAcceptConsentRequest()
	acc.SetGrantScope(req.GrantScope)
	acc.SetSession(*session)
	acc.SetGrantAccessTokenAudience(hreq.RequestedAccessTokenAudience)
	acc.SetRemember(true)
	acc.SetRememberFor(3600)

	accReq, _, err := a.hclient.AdminApi.AcceptConsentRequest(r.Context()).ConsentChallenge(req.Challenge).AcceptConsentRequest(*acc).Execute()
	if err != nil {
		return err
	}
	resp.Redirect = accReq.RedirectTo
	if err := a.resEncoder(w, r, &resp); err != nil {
		return err
	}
	return nil
}
