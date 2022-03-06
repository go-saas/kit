package service

import (
	"github.com/goxiaoy/go-saas-kit/pkg/authn"
	v1 "github.com/goxiaoy/go-saas-kit/saas/api/tenant/v1"
	"github.com/goxiaoy/go-saas/common"
)

type UserTenantContributor struct {
	us *UserService
}

var _ common.TenantResolveContributor = (*UserTenantContributor)(nil)

func NewUserTenantContributor(us *UserService) *UserTenantContributor {
	return &UserTenantContributor{us: us}
}

func (u *UserTenantContributor) Name() string {
	return "UserTenant"
}

func (u *UserTenantContributor) Resolve(trCtx *common.TenantResolveContext) error {
	ui, _ := authn.FromUserContext(trCtx.Context())
	if len(ui.GetId()) > 0 {
		//user logged in
		ok, err := u.us.CheckUserTenantInternal(trCtx.Context(), ui.GetId(), trCtx.TenantIdOrName)
		if err != nil {
			return err
		}
		if !ok {
			return v1.ErrorTenantForbidden("")
		}
	}
	return nil
}
