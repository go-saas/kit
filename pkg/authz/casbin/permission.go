package casbin

import (
	"context"
	"errors"
	"fmt"
	"github.com/casbin/casbin/v2/util"
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authorization"
	"github.com/goxiaoy/go-saas/common"
)

type PermissionService struct {
	enforcer *EnforcerProvider
}

func NewPermissionService(enforcer *EnforcerProvider) *PermissionService {
	return &PermissionService{
		enforcer: enforcer,
	}
}

var _ authorization.PermissionManagementService = (*PermissionService)(nil)
var _ authorization.PermissionChecker = (*PermissionService)(nil)

func (p *PermissionService) IsGrant(ctx context.Context, resource authorization.Resource, action authorization.Action, subjects ...authorization.Subject) (authorization.Effect, error) {
	tenantInfo := common.FromCurrentTenant(ctx)
	return p.IsGrantTenant(ctx, resource, action, tenantInfo.GetId(), subjects...)
}

func (p *PermissionService) IsGrantTenant(ctx context.Context, resource authorization.Resource, action authorization.Action, tenantID string, subjects ...authorization.Subject) (authorization.Effect, error) {
	enforcer, err := p.enforcer.Get(ctx)
	if err != nil {
		return authorization.EffectForbidden, err
	}
	subs := make([][]interface{}, len(subjects))
	for i, subject := range subjects {
		subs[i] = []interface{}{subject.GetIdentity(), resource.GetNamespace(), resource.GetIdentity(), action.GetIdentity(), tenantID}
	}
	results, err := enforcer.BatchEnforce(subs)
	if err != nil {
		return authorization.EffectForbidden, err
	}
	var grant bool
	for i := range results {
		if results[i] {
			grant = true
		}
	}
	if grant {
		return authorization.EffectGrant, nil
	}
	return authorization.EffectForbidden, nil
}

func (p *PermissionService) AddGrant(ctx context.Context, resource authorization.Resource, action authorization.Action, subject authorization.Subject, tenantID string, effect authorization.Effect) error {
	enforcer, err := p.enforcer.Get(ctx)
	if err != nil {
		return err
	}

	eff, err := mapToEffect(effect)
	if err != nil {
		return err
	}
	_, err = enforcer.AddPolicy(subject.GetIdentity(), resource.GetNamespace(), resource.GetIdentity(), action.GetIdentity(), tenantID, eff)
	if err != nil {
		return err
	}
	return nil
}

func (p *PermissionService) ListAcl(ctx context.Context, subjects ...authorization.Subject) ([]authorization.PermissionBean, error) {
	//list
	enforcer, err := p.enforcer.Get(ctx)
	if err != nil {
		return nil, err
	}
	policies := enforcer.GetPolicy()
	var ret []authorization.PermissionBean
	for _, policy := range policies {
		for _, subject := range subjects {
			if util.KeyMatch(policy[0], subject.GetIdentity()) {
				ret = append(ret, authorization.NewPermissionBean(authorization.NewEntityResource(policy[1], policy[2]),
					authorization.ActionStr(policy[3]),
					authorization.SubjectStr(policy[0]),
					policy[4], mapToAuthEffect(policy[5]),
				))
			}
		}

	}
	return ret, nil
}

func (p *PermissionService) UpdateGrant(ctx context.Context, subject authorization.Subject, acl []authorization.UpdateSubjectPermission) error {
	enforcer, err := p.enforcer.Get(ctx)
	if err != nil {
		return err
	}
	_, err = enforcer.RemoveFilteredPolicy(0, subject.GetIdentity())
	if err != nil {
		return err
	}

	rules := make([][]string, len(acl))
	for i, permission := range acl {
		eff, err := mapToEffect(permission.Effect)
		if err != nil {
			return err
		}
		rules[i] = []string{subject.GetIdentity(),
			permission.Resource.GetNamespace(),
			permission.Resource.GetIdentity(),
			permission.Action.GetIdentity(),
			permission.TenantID,
			eff}
	}
	_, err = enforcer.AddPolicies(rules)
	if err != nil {
		return err
	}
	return nil
}

func mapToEffect(effect authorization.Effect) (string, error) {
	eff := "allow"
	if effect == authorization.EffectGrant {
		eff = "allow"
	} else if effect == authorization.EffectForbidden {
		eff = "deny"
	} else {
		return "", errors.New(fmt.Sprintf("effect should be one of %s,%s", "grant", "forbidden"))
	}
	return eff, nil
}

func mapToAuthEffect(eff string) authorization.Effect {
	if eff == "allow" {
		return authorization.EffectGrant
	} else if eff == "deny" {
		return authorization.EffectForbidden
	}
	return authorization.EffectUnknown
}

var PermissionProviderSet = wire.NewSet(
	NewPermissionService,
	wire.Bind(new(authorization.PermissionManagementService), new(*PermissionService)),
	wire.Bind(new(authorization.PermissionChecker), new(*PermissionService)),
)
