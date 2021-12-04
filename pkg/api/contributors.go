package api

import (
	"context"
	"github.com/goxiaoy/go-saas-kit/auth/current"
	"github.com/goxiaoy/go-saas/common"
	shttp "github.com/goxiaoy/go-saas/common/http"
)

type SaasContributor struct {
	hmtOpt *shttp.WebMultiTenancyOption
}

var _ Contributor = (*SaasContributor)(nil)

func NewSaasContributor(hmtOpt *shttp.WebMultiTenancyOption) *SaasContributor {
	return &SaasContributor{
		hmtOpt: hmtOpt,
	}
}

func (s *SaasContributor) RecoverContext(ctx context.Context, headers Header) (context.Context, error) {
	if !headers.HasKey(s.hmtOpt.TenantKey) {
		return ctx, nil
	}
	tenantId := headers.Get(s.hmtOpt.TenantKey)
	return common.NewCurrentTenant(ctx, tenantId, ""), nil
}

func (s *SaasContributor) CreateHeader(ctx context.Context) map[string]string {
	ti := common.FromCurrentTenant(ctx)
	return map[string]string{
		s.hmtOpt.TenantKey: ti.Id,
	}
}

type UserContributor struct {
}

var _ Contributor = (*UserContributor)(nil)

func NewUserContributor() *UserContributor {
	return &UserContributor{}
}

func (u *UserContributor) RecoverContext(ctx context.Context, headers Header) (context.Context, error) {
	if !headers.HasKey("user") {
		return ctx, nil
	}
	user := headers.Get("user")
	return current.NewUserContext(ctx, current.NewUserInfo(user)), nil
}

func (u *UserContributor) CreateHeader(ctx context.Context) map[string]string {
	res := map[string]string{}
	if userInfo, ok := current.FromUserContext(ctx); ok {
		res["user"] = userInfo.GetId()
	}
	return res
}
