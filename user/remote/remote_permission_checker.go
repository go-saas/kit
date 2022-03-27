package remote

import (
	"context"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authz"
	v1 "github.com/goxiaoy/go-saas-kit/user/api/permission/v1"
	"github.com/goxiaoy/go-saas-kit/user/util"
	"github.com/samber/lo"
)

type PermissionChecker struct {
	client v1.PermissionServiceClient
}

var _ authz.PermissionChecker = (*PermissionChecker)(nil)
var _ authz.PermissionManagementService = (*PermissionChecker)(nil)

func NewRemotePermissionChecker(client v1.PermissionServiceClient) *PermissionChecker {
	return &PermissionChecker{
		client: client,
	}
}

func (r *PermissionChecker) IsGrantTenant(ctx context.Context, resource authz.Resource, action authz.Action, tenantID string, subjects ...authz.Subject) (authz.Effect, error) {
	var protoSubs = make([]string, len(subjects))
	for i, subject := range subjects {
		protoSubs[i] = subject.GetIdentity()
	}
	grant, err := r.client.CheckForSubjects(ctx, &v1.CheckSubjectsPermissionRequest{
		Namespace: resource.GetNamespace(),
		Resource:  resource.GetIdentity(),
		Action:    action.GetIdentity(),
		Subjects:  protoSubs,
		TenantId:  tenantID,
	})
	if err != nil {
		return authz.EffectForbidden, err
	}
	if grant.Effect == v1.Effect_GRANT {
		return authz.EffectGrant, nil
	}
	return authz.EffectForbidden, nil
}

func (r *PermissionChecker) AddGrant(ctx context.Context, resource authz.Resource, action authz.Action, subject authz.Subject, tenantID string, effect authz.Effect) error {
	_, err := r.client.AddSubjectPermission(ctx, &v1.AddSubjectPermissionRequest{
		Namespace: resource.GetNamespace(),
		Resource:  resource.GetIdentity(),
		Action:    action.GetIdentity(),
		Subject:   subject.GetIdentity(),
		Effect:    util.MapAuthEffect2PbEffect(effect),
		TenantId:  tenantID,
	})
	return err
}

func (r *PermissionChecker) ListAcl(ctx context.Context, subjects ...authz.Subject) ([]authz.PermissionBean, error) {
	subs := make([]string, len(subjects))
	for i, subject := range subjects {
		subs[i] = subject.GetIdentity()
	}
	acl, err := r.client.ListSubjectPermission(ctx, &v1.ListSubjectPermissionRequest{Subjects: subs})
	if err != nil {
		return nil, err
	}
	res := make([]authz.PermissionBean, len(acl.Acl))
	for i, permission := range acl.Acl {
		j := authz.PermissionBean{}
		util.MapPbPermissionToBean(permission, &j)
		res[i] = j
	}
	return res, nil
}

func (r *PermissionChecker) UpdateGrant(ctx context.Context, subject authz.Subject, acl []authz.UpdateSubjectPermission) error {
	var pbAcl = lo.Map(acl, func(a authz.UpdateSubjectPermission, _ int) *v1.UpdateSubjectPermissionAcl {
		return &v1.UpdateSubjectPermissionAcl{
			Namespace: a.Resource.GetNamespace(),
			Resource:  a.Resource.GetIdentity(),
			Action:    a.Action.GetIdentity(),
			Effect:    util.MapAuthEffect2PbEffect(a.Effect),
			TenantId:  a.TenantID,
		}
	})
	_, err := r.client.UpdateSubjectPermission(ctx, &v1.UpdateSubjectPermissionRequest{
		Subject: subject.GetIdentity(),
		Acl:     pbAcl,
	})
	return err
}

func (r *PermissionChecker) RemoveGrant(ctx context.Context, resource authz.Resource, action authz.Action, subject authz.Subject, tenantID string, effects []authz.Effect) error {
	var effs = lo.Map(effects, func(e authz.Effect, _ int) v1.Effect {
		return util.MapAuthEffect2PbEffect(e)
	})
	_, err := r.client.RemoveSubjectPermission(ctx, &v1.RemoveSubjectPermissionRequest{
		Namespace: resource.GetNamespace(),
		Resource:  resource.GetIdentity(),
		Action:    action.GetIdentity(),
		Subject:   subject.GetIdentity(),
		Effects:   effs,
		TenantId:  tenantID,
	})
	return err
}
