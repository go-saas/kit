package http

import (
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	http2 "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/kratos/v2/transport/http/binding"
	"github.com/go-saas/kit/oidc/service"
	"github.com/go-saas/kit/pkg/authn"
	"github.com/go-saas/kit/pkg/utils"
	v1 "github.com/go-saas/kit/user/api/auth/v1"
	"github.com/go-saas/kit/user/private/biz"
	client "github.com/ory/hydra-client-go/v2"
	"net/http"
)

type Auth struct {
	reqDecoder http2.DecodeRequestFunc
	um         *biz.UserManager
	logger     *log.Helper
	signIn     *biz.SignInManager
	hclient    *client.APIClient
}

func NewAuth(
	reqDecoder http2.DecodeRequestFunc,
	um *biz.UserManager,
	l log.Logger,
	signIn *biz.SignInManager,
	hclient *client.APIClient,
) *Auth {
	return &Auth{
		reqDecoder: reqDecoder,
		um:         um,
		logger:     log.NewHelper(l),
		signIn:     signIn,
		hclient:    hclient,
	}
}

func (a *Auth) LoginGet(w http.ResponseWriter, r *http.Request) (*v1.GetLoginResponse, error) {
	var req v1.GetLoginRequest
	if err := binding.BindQuery(r.URL.Query(), &req); err != nil {
		return nil, err
	}
	var resp = &v1.GetLoginResponse{}
	ui, _ := authn.FromUserContext(r.Context())

	if len(ui.GetId()) > 0 {
		resp.Redirect = req.Redirect
		if len(resp.Redirect) == 0 {
			resp.Redirect = "/"
		}
	}

	if len(req.LoginChallenge) > 0 {
		//hydra
		if hreq, raw, err := a.hclient.OAuth2Api.GetOAuth2LoginRequest(r.Context()).LoginChallenge(req.LoginChallenge).Execute(); err != nil {
			return resp, service.TransformHydraErr(raw, err)
		} else {
			// If hydra was already able to authenticate the user, skip will be true and we do not need to re-authenticate
			// the user.
			if hreq.Skip {
				acc, raw, err := a.hclient.OAuth2Api.AcceptOAuth2LoginRequest(r.Context()).
					LoginChallenge(hreq.Challenge).
					AcceptOAuth2LoginRequest(*client.NewAcceptOAuth2LoginRequest(hreq.Subject)).
					Execute()
				if err != nil {
					return resp, service.TransformHydraErr(raw, err)
				}
				resp.Redirect = acc.RedirectTo
			} else if len(ui.GetId()) > 0 {
				//user logged in
				acc, raw, err := a.hclient.OAuth2Api.AcceptOAuth2LoginRequest(r.Context()).
					LoginChallenge(hreq.Challenge).
					AcceptOAuth2LoginRequest(*client.NewAcceptOAuth2LoginRequest(ui.GetId())).
					Execute()
				if err != nil {
					return resp, service.TransformHydraErr(raw, err)
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
	return resp, nil
}

func (a *Auth) LoginPost(w http.ResponseWriter, r *http.Request) (*v1.WebLoginAuthReply, error) {
	var req v1.WebLoginAuthRequest
	if err := a.reqDecoder(r, &req); err != nil {
		return nil, err
	}
	var resp = &v1.WebLoginAuthReply{}
	if len(req.Challenge) > 0 {
		// Let's see if the user decided to accept or reject the consent request..
		if req.Reject {
			// Looks like the consent request was denied by the user
			reject := *client.NewRejectOAuth2Request()
			//TODO error
			reject.SetError("access_denied")
			reject.SetErrorDescription("The resource owner denied the request")
			hreq, raw, err := a.hclient.OAuth2Api.RejectOAuth2LoginRequest(r.Context()).LoginChallenge(req.Challenge).RejectOAuth2Request(reject).Execute()
			if err != nil {
				return resp, service.TransformHydraErr(raw, err)
			}
			resp.Redirect = hreq.RedirectTo
			//return
			return resp, nil
		}
	}
	//validate sign in
	err, uid := a.signIn.PasswordSignInWithUsername(r.Context(), req.Username, req.Password, req.Remember, true)
	if err != nil {
		return resp, err
	}
	if len(req.Challenge) > 0 {
		// Seems like the user authenticated! Let's tell hydra...
		acc := *client.NewAcceptOAuth2LoginRequest(uid)
		acc.SetRemember(req.Remember)
		//TODO from config
		acc.SetRememberFor(3600)
		acc.SetSubject(uid)
		hreq, raw, err := a.hclient.OAuth2Api.AcceptOAuth2LoginRequest(r.Context()).
			LoginChallenge(req.Challenge).
			AcceptOAuth2LoginRequest(acc).Execute()
		if err != nil {
			return resp, service.TransformHydraErr(raw, err)
		}
		resp.Redirect = hreq.RedirectTo
	}
	return resp, nil
}

func (a *Auth) LoginOutGet(w http.ResponseWriter, r *http.Request) (*v1.GetLogoutResponse, error) {
	var req v1.GetLogoutRequest
	if err := binding.BindQuery(r.URL.Query(), &req); err != nil {
		return nil, err
	}
	var resp = &v1.GetLogoutResponse{}
	if len(req.LogoutChallenge) > 0 {
		_, raw, err := a.hclient.OAuth2Api.GetOAuth2LogoutRequest(r.Context()).LogoutChallenge(req.LogoutChallenge).Execute()
		if err != nil {
			return resp, service.TransformHydraErr(raw, err)
		}
		resp.Challenge = req.LogoutChallenge
	}

	return resp, nil
}

func (a *Auth) Logout(w http.ResponseWriter, r *http.Request) (*v1.LogoutResponse, error) {
	var req v1.LogoutRequest
	if err := a.reqDecoder(r, &req); err != nil {
		return nil, err
	}
	var resp = &v1.LogoutResponse{}
	err := a.signIn.SignOut(r.Context())
	if err != nil {
		return resp, err
	}
	if len(req.Challenge) > 0 {
		hreq, raw, err := a.hclient.OAuth2Api.AcceptOAuth2LogoutRequest(r.Context()).LogoutChallenge(req.Challenge).Execute()
		if err != nil {
			return resp, service.TransformHydraErr(raw, err)
		}
		resp.Redirect = hreq.RedirectTo
	}

	return resp, nil

}

func (a *Auth) ConsentGet(w http.ResponseWriter, r *http.Request) (*v1.GetConsentResponse, error) {
	ui, err := authn.ErrIfUnauthenticated(r.Context())
	if err != nil {
		return nil, err
	}
	var req v1.GetConsentRequest
	if err := binding.BindQuery(r.URL.Query(), &req); err != nil {
		return nil, err
	}
	var resp = &v1.GetConsentResponse{}
	if len(req.ConsentChallenge) == 0 {
		return resp, errors.BadRequest("CONSENT_CHALLENGE_REQUIRED", "")
	}
	hreq, raw, err := a.hclient.OAuth2Api.GetOAuth2ConsentRequest(r.Context()).ConsentChallenge(req.ConsentChallenge).Execute()
	if err != nil {
		return resp, service.TransformHydraErr(raw, err)
	}
	if ui.GetId() != hreq.GetSubject() {
		return nil, errors.Unauthorized("SUB_MISMATCH", "")
	}
	if hreq.GetSkip() {
		acc := *client.NewAcceptOAuth2ConsentRequest()
		acc.SetGrantScope(hreq.RequestedScope)
		acc.SetGrantAccessTokenAudience(hreq.RequestedAccessTokenAudience)
		//acc.SetSession(client.ConsentRequestSession{
		//	AccessToken: nil,
		//	IdToken:     nil,
		//})
		accReq, raw, err := a.hclient.OAuth2Api.AcceptOAuth2ConsentRequest(r.Context()).ConsentChallenge(req.ConsentChallenge).AcceptOAuth2ConsentRequest(acc).Execute()
		if err != nil {
			return resp, service.TransformHydraErr(raw, err)
		}
		resp.Redirect = accReq.RedirectTo
		return resp, nil
	}
	resp.Challenge = hreq.Challenge
	resp.RequestedScope = hreq.RequestedScope
	resp.UserId = hreq.GetSubject()
	resp.Client = mapClients(hreq.GetClient())

	return resp, nil
}

func (a *Auth) Consent(w http.ResponseWriter, r *http.Request) (*v1.GrantConsentResponse, error) {
	userInfo, err := authn.ErrIfUnauthenticated(r.Context())
	if err != nil {
		return nil, err
	}
	var req v1.GrantConsentRequest
	if err := a.reqDecoder(r, &req); err != nil {
		return nil, err
	}
	if len(req.Challenge) == 0 {
		return nil, errors.BadRequest("CHALLENGE_REQUIRED", "")
	}

	var resp = &v1.GrantConsentResponse{}

	if req.Reject {
		reject := *client.NewRejectOAuth2Request()
		//TODO
		reject.SetError("access_denied")
		reject.SetErrorDescription("The resource owner denied the request")
		hreq, raw, err := a.hclient.OAuth2Api.RejectOAuth2ConsentRequest(r.Context()).ConsentChallenge(req.Challenge).RejectOAuth2Request(reject).Execute()
		if err != nil {
			return resp, service.TransformHydraErr(raw, err)
		}
		resp.Redirect = hreq.RedirectTo
		//return
		return resp, nil
	}
	//user allow
	// The session allows us to set session data for id and access tokens
	session := client.NewAcceptOAuth2ConsentRequestSession()
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

	hreq, raw, err := a.hclient.OAuth2Api.GetOAuth2ConsentRequest(r.Context()).ConsentChallenge(req.Challenge).Execute()
	if err != nil {
		return resp, service.TransformHydraErr(raw, err)
	}

	if hreq.Subject == nil || *hreq.Subject != userInfo.GetId() {
		return nil, errors.Unauthorized("", "")
	}

	acc := client.NewAcceptOAuth2ConsentRequest()
	acc.SetGrantScope(req.GrantScope)
	acc.SetSession(*session)
	acc.SetGrantAccessTokenAudience(hreq.RequestedAccessTokenAudience)
	acc.SetRemember(true)
	acc.SetRememberFor(3600)

	accReq, raw, err := a.hclient.OAuth2Api.AcceptOAuth2ConsentRequest(r.Context()).ConsentChallenge(req.Challenge).AcceptOAuth2ConsentRequest(*acc).Execute()
	if err != nil {
		return resp, service.TransformHydraErr(raw, err)
	}
	resp.Redirect = accReq.RedirectTo
	return resp, nil
}

func mapClients(c client.OAuth2Client) *v1.OAuthClient {
	ret := &v1.OAuthClient{
		ClientId:   c.ClientId,
		ClientName: c.ClientName,

		ClientUri: c.ClientUri,
		Contacts:  c.Contacts,

		LogoUri:   c.LogoUri,
		Metadata:  utils.Map2Structpb(c.Metadata.(map[string]interface{})),
		Owner:     c.Owner,
		PolicyUri: c.PolicyUri,
	}
	return ret
}
