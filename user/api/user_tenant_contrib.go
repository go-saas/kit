package api

import (
	"github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/authn"
	v12 "github.com/go-saas/kit/saas/api/tenant/v1"
	v1 "github.com/go-saas/kit/user/api/user/v1"
	"github.com/go-saas/saas"
)

// UserTenantContrib impl saas.TenantResolveContrib from calling remote or local service.
//
// check whether user can present in a tenant
type UserTenantContrib struct {
	srv v1.UserInternalServiceServer
}

func NewUserTenantContrib(client v1.UserInternalServiceServer) *UserTenantContrib {
	return &UserTenantContrib{
		srv: client,
	}
}

var _ saas.TenantResolveContrib = (*UserTenantContrib)(nil)

func (u *UserTenantContrib) Name() string {
	return "UserTenant"
}

func (u *UserTenantContrib) Resolve(trCtx *saas.Context) error {
	ui, _ := authn.FromUserContext(trCtx.Context())
	if len(ui.GetId()) > 0 { //user logged in
		//replace withe trusted environment to skip trusted check if in same process
		trCtx.WithContext(api.NewTrustedContext(trCtx.Context()))
		ok, err := u.srv.CheckUserTenant(trCtx.Context(), &v1.CheckUserTenantRequest{
			UserId:   ui.GetId(),
			TenantId: trCtx.TenantIdOrName,
		})
		if err != nil {
			return err
		}
		if !ok.Ok {
			return v12.ErrorTenantForbiddenLocalized(trCtx.Context(), nil, nil)
		}
	}
	return nil
}
