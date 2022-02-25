package authz

import (
	"context"
	"fmt"
	"github.com/goxiaoy/go-saas/common"
)

type PermissionChecker interface {
	IsGrant(ctx context.Context, resource Resource, action Action, subjects ...Subject) (Effect, error)
	IsGrantTenant(ctx context.Context, resource Resource, action Action, tenantID string, subjects ...Subject) (Effect, error)
}

var _ PermissionChecker = (*PermissionService)(nil)

func (p *PermissionService) IsGrant(ctx context.Context, resource Resource, action Action, subjects ...Subject) (Effect, error) {
	tenantInfo, _ := common.FromCurrentTenant(ctx)
	return p.IsGrantTenant(ctx, resource, action, tenantInfo.GetId(), subjects...)
}

func (p *PermissionService) IsGrantTenant(ctx context.Context, resource Resource, action Action, tenantID string, subjects ...Subject) (Effect, error) {
	p.mux.Lock()
	defer p.mux.Unlock()
	var anyAllow bool
	//TODO host side?

	for _, subject := range subjects {
		for _, bean := range p.v {
			if match(bean.Namespace, resource.GetNamespace()) &&
				match(bean.Subject, subject.GetIdentity()) &&
				match(bean.Resource, resource.GetIdentity()) &&
				match(bean.Action, action.GetIdentity()) &&
				match(bean.TenantID, tenantID) {
				if bean.Effect == EffectForbidden {
					p.log.Debugf("Subject %s Action %s to Resource %s forbidden", subject.GetIdentity(), action.GetIdentity(), fmt.Sprintf("%s/%s", resource.GetNamespace(), resource.GetIdentity()))
					return EffectForbidden, nil
				}
				if bean.Effect == EffectGrant {
					anyAllow = true
				}
			}
		}
		if anyAllow {
			p.log.Debugf("Subject %s Action %s to Resource %s grant", subject.GetIdentity(), action.GetIdentity(), fmt.Sprintf("%s/%s", resource.GetNamespace(), resource.GetIdentity()))
			return EffectGrant, nil
		}
		p.log.Debugf("Subject %s Action %s to Resource %s unknown", subject.GetIdentity(), action.GetIdentity(), fmt.Sprintf("%s/%s", resource.GetNamespace(), resource.GetIdentity()))
	}
	return EffectUnknown, nil
}
