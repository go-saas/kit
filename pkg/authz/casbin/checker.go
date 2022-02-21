package casbin

import (
	"context"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authz"
	"github.com/goxiaoy/go-saas/common"
)

var _ authz.PermissionChecker = (*PermissionService)(nil)

func (p *PermissionService) IsGrant(ctx context.Context, resource authz.Resource, action authz.Action, subjects ...authz.Subject) (authz.Effect, error) {
	tenantInfo := common.FromCurrentTenant(ctx)
	return p.IsGrantTenant(ctx, resource, action, tenantInfo.GetId(), subjects...)
}

func (p *PermissionService) IsGrantTenant(ctx context.Context, resource authz.Resource, action authz.Action, tenantID string, subjects ...authz.Subject) (authz.Effect, error) {
	enforcer, err := p.enforcer.Get(ctx)
	if err != nil {
		return authz.EffectForbidden, err
	}
	subs := make([][]interface{}, len(subjects))
	for i, subject := range subjects {
		subs[i] = []interface{}{subject.GetIdentity(), resource.GetNamespace(), resource.GetIdentity(), action.GetIdentity(), tenantID}
	}
	results, err := enforcer.BatchEnforce(subs)
	if err != nil {
		return authz.EffectForbidden, err
	}
	var grant bool
	for i := range results {
		if results[i] {
			grant = true
		}
	}
	if grant {
		return authz.EffectGrant, nil
	}
	return authz.EffectForbidden, nil
}
