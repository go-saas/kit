package remote

import (
	"github.com/goxiaoy/go-saas-kit/pkg/authn"
	v12 "github.com/goxiaoy/go-saas-kit/saas/api/tenant/v1"
	v1 "github.com/goxiaoy/go-saas-kit/user/api/user/v1"
	"github.com/goxiaoy/go-saas/common"
)

type UserTenantContributor struct {
	client v1.UserServiceClient
}

func NewUserTenantContributor(client v1.UserServiceClient) *UserTenantContributor {
	return &UserTenantContributor{
		client: client,
	}
}

var _ common.TenantResolveContributor = (*UserTenantContributor)(nil)

func (u *UserTenantContributor) Name() string {
	return "RemoteUserTenant"
}

func (u *UserTenantContributor) Resolve(trCtx *common.TenantResolveContext) error {
	ui, _ := authn.FromUserContext(trCtx.Context())
	if len(ui.GetId()) > 0 {
		//user logged in
		ok, err := u.client.CheckUserTenant(trCtx.Context(), &v1.CheckUserTenantRequest{
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
