package authz

import (
	"context"
	"github.com/go-saas/saas"
)

type PermissionManagementService interface {
	AddGrant(ctx context.Context, resource Resource, action Action, subject Subject, tenantID string, effect Effect) error
	//ListAcl list permission of subjects. if not subjects provided, all acl will be returned
	ListAcl(ctx context.Context, subjects ...Subject) ([]PermissionBean, error)
	UpdateGrant(ctx context.Context, subject Subject, acl []UpdateSubjectPermission) error
	RemoveGrant(ctx context.Context, resource Resource, action Action, subject Subject, tenantID string, effects []Effect) error
}

func EnsureGrant(ctx context.Context, mgr PermissionManagementService, checker PermissionChecker, resource Resource, action Action, subject Subject, tenantID string) error {
	eff, err := checker.IsGrantTenant(ctx, []*Requirement{NewRequirement(resource, action)}, tenantID, subject)
	if err != nil {
		return err
	}
	if eff[0] != EffectGrant {
		err = mgr.AddGrant(ctx, resource, action, subject, tenantID, EffectGrant)
		if err != nil {
			return err
		}
	}
	return nil
}

func EnsureForbidden(ctx context.Context, mgr PermissionManagementService, checker PermissionChecker, resource Resource, action Action, subject Subject, tenantID string) error {
	eff, err := checker.IsGrantTenant(ctx, []*Requirement{NewRequirement(resource, action)}, tenantID, subject)
	if err != nil {
		return err
	}
	if eff[0] != EffectForbidden {
		err = mgr.AddGrant(ctx, resource, action, subject, tenantID, EffectForbidden)
		if err != nil {
			return err
		}
	}
	return nil
}

func NormalizeTenantId(ctx context.Context, tenantId string) string {
	ti, _ := saas.FromCurrentTenant(ctx)
	if ti.GetId() == "" {
		//host side
		return tenantId
	}
	return ti.GetId()
}
