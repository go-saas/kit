package api

import (
	"github.com/goxiaoy/go-saas"
	"github.com/goxiaoy/go-saas-kit/pkg/api"
	"github.com/goxiaoy/go-saas-kit/pkg/authn"
	v12 "github.com/goxiaoy/go-saas-kit/saas/api/tenant/v1"
	v1 "github.com/goxiaoy/go-saas-kit/user/api/user/v1"
)

// UserTenantContrib impl saas.TenantResolveContrib from calling remote or local service.
//
// check whether user can present in a tenant
type UserTenantContrib struct {
	srv v1.UserServiceServer
}

func NewUserTenantContrib(client v1.UserServiceServer) *UserTenantContrib {
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
			return v12.ErrorTenantForbidden("")
		}
	}
	return nil
}