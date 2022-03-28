package casbin

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authz"
)

var _ authz.PermissionChecker = (*PermissionService)(nil)

func (p *PermissionService) IsGrantTenant(ctx context.Context, requirements authz.RequirementList, tenantID string, subjects ...authz.Subject) ([]authz.Effect, error) {
	res := []authz.Effect{}
	for _, requirement := range requirements {
		eff, err := p.isGrantTenant(ctx, requirement.Resource, requirement.Action, tenantID, subjects...)
		if err != nil {
			return nil, err
		}
		res = append(res, eff)
	}
	return res, nil
}

func (p *PermissionService) isGrantTenant(ctx context.Context, resource authz.Resource, action authz.Action, tenantID string, subjects ...authz.Subject) (authz.Effect, error) {
	//find permission definition of current resource and action

	def := authz.MustFindDef(resource.GetNamespace(), action)

	if (def.Side == authz.PermissionHostSideOnly && len(tenantID) != 0) || (def.Side == authz.PermissionTenantSideOnly && len(tenantID) == 0) {
		return authz.EffectForbidden, errors.New(400, authz.DefNotFoundReason, fmt.Sprintf("action %s in %s side mismatch",
			action.GetIdentity(), resource.GetNamespace()))
	}
	if def.IsInternalOnly() {
		//internal ignore tenant
		tenantID = "*"
	}

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
