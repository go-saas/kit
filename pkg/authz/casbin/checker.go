package casbin

import (
	"context"
	"github.com/casbin/casbin/v2"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-saas/kit/pkg/authz/authz"
)

var _ authz.PermissionChecker = (*PermissionService)(nil)

func (p *PermissionService) IsGrantTenant(ctx context.Context, requirements authz.RequirementList, tenantID string, subjects ...authz.Subject) ([]authz.Effect, error) {
	var res []authz.Effect

	enforcer, err := p.enforcer.Get(ctx)
	if err != nil {
		return nil, err
	}

	for _, requirement := range requirements {
		eff, err := p.isGrantTenant(ctx, enforcer, requirement.Resource, requirement.Action, tenantID, subjects...)
		if err != nil {
			return nil, err
		}
		res = append(res, eff)
	}
	return res, nil
}

func (p *PermissionService) isGrantTenant(ctx context.Context, enforcer *casbin.SyncedEnforcer, resource authz.Resource, action authz.Action, tenantID string, subjects ...authz.Subject) (authz.Effect, error) {
	//find permission definition of current resource and action

	def, err := authz.FindDef(resource.GetNamespace(), action, false)
	if err != nil {
		if errors.Reason(err) == authz.DefNotFoundReason {
			//just forbid
			return authz.EffectForbidden, nil
		} else {
			return authz.EffectForbidden, err
		}
	}

	if (def.Side == authz.PermissionAllowSide_HOST_ONLY && len(tenantID) != 0) || (def.Side == authz.PermissionAllowSide_TENANT_ONLY && len(tenantID) == 0) {
		//just forbid
		return authz.EffectForbidden, nil
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
