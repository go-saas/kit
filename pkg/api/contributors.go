package api

import (
	"context"
	"github.com/goxiaoy/go-saas-kit/pkg/authn"
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
	tenantId := headers.Get(s.hmtOpt.TenantKey)
	return common.NewCurrentTenant(ctx, tenantId, ""), nil
}

func (s *SaasContributor) CreateHeader(ctx context.Context) map[string]string {
	ti := common.FromCurrentTenant(ctx)
	return map[string]string{
		s.hmtOpt.TenantKey: ti.GetId(),
	}
}

type UserContributor struct {
}

var _ Contributor = (*UserContributor)(nil)

func NewUserContributor() *UserContributor {
	return &UserContributor{}
}

func (u *UserContributor) RecoverContext(ctx context.Context, headers Header) (context.Context, error) {
	user := headers.Get("user")
	return authn.NewUserContext(ctx, authn.NewUserInfo(user)), nil
}

func (u *UserContributor) CreateHeader(ctx context.Context) map[string]string {
	res := map[string]string{}
	if userInfo, ok := authn.FromUserContext(ctx); ok {
		res["user"] = userInfo.GetId()
	}
	return res
}

type ClientContributor struct {
}

var _ Contributor = (*ClientContributor)(nil)

func NewClientContributor() *ClientContributor {
	return &ClientContributor{}
}

func (u *ClientContributor) RecoverContext(ctx context.Context, headers Header) (context.Context, error) {
	client := headers.Get("client")
	return authn.NewClientContext(ctx, client), nil
}

func (u *ClientContributor) CreateHeader(ctx context.Context) map[string]string {
	res := map[string]string{}
	if client, ok := authn.FromClientContext(ctx); ok {
		res["client"] = client
	}
	return res
}
