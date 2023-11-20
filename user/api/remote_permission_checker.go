package api

import (
	"context"
	"github.com/go-saas/kit/pkg/authz/authz"
	v1 "github.com/go-saas/kit/user/api/permission/v1"
	"github.com/go-saas/kit/user/util"
	"github.com/samber/lo"
)

// PermissionChecker impl authz.PermissionChecker and authz.PermissionManagementService from calling remote service
type PermissionChecker struct {
	srv v1.PermissionInternalServiceServer
}

var _ authz.PermissionChecker = (*PermissionChecker)(nil)
var _ authz.PermissionManagementService = (*PermissionChecker)(nil)

func NewRemotePermissionChecker(srv v1.PermissionInternalServiceServer) *PermissionChecker {
	return &PermissionChecker{
		srv: srv,
	}
}

func (r *PermissionChecker) IsGrantTenant(ctx context.Context, requirements authz.RequirementList, tenantID string, subjects ...authz.Subject) ([]authz.Effect, error) {
	var protoSubs = make([]string, len(subjects))
	for i, subject := range subjects {
		protoSubs[i] = subject.GetIdentity()
	}
	grantResp, err := r.srv.CheckForSubjects(ctx, &v1.CheckSubjectsPermissionRequest{
		Requirements: lo.Map(requirements, func(t *authz.Requirement, _ int) *v1.PermissionRequirement {
			return &v1.PermissionRequirement{
				Namespace: t.Resource.GetNamespace(),
				Resource:  t.Resource.GetIdentity(),
				Action:    t.Action.GetIdentity(),
			}
		}),
		Subjects: protoSubs,
		TenantId: tenantID,
	})
	if err != nil {
		return nil, err
	}
	effList := lo.Map(grantResp.EffectList, func(eff v1.Effect, _ int) authz.Effect {
		switch eff {
		case v1.Effect_GRANT:
			return authz.EffectGrant
		case v1.Effect_FORBIDDEN:
			return authz.EffectForbidden
		case v1.Effect_UNKNOWN:
			return authz.EffectUnknown
		}
		return authz.EffectUnknown
	})
	return effList, nil
}

func (r *PermissionChecker) AddGrant(ctx context.Context, resource authz.Resource, action authz.Action, subject authz.Subject, tenantID string, effect authz.Effect) error {
	_, err := r.srv.AddSubjectPermission(ctx, &v1.AddSubjectPermissionRequest{
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
	acl, err := r.srv.ListSubjectPermission(ctx, &v1.ListSubjectPermissionRequest{Subjects: subs})
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
	_, err := r.srv.UpdateSubjectPermission(ctx, &v1.UpdateSubjectPermissionRequest{
		Subject: subject.GetIdentity(),
		Acl:     pbAcl,
	})
	return err
}

func (r *PermissionChecker) RemoveGrant(ctx context.Context, subject authz.Subject, filter ...authz.FilterFunc) error {
	f := &authz.Filter{}
	for _, filterFunc := range filter {
		filterFunc(f)
	}

	var effs = lo.Map(f.Effects, func(e authz.Effect, _ int) v1.Effect {
		return util.MapAuthEffect2PbEffect(e)
	})
	req := &v1.RemoveSubjectPermissionRequest{
		Subject:  subject.GetIdentity(),
		Effects:  effs,
		TenantId: f.TenantID,
	}
	if f.Resource != nil {
		n := f.Resource.GetNamespace()
		req.Namespace = &n
		i := f.Resource.GetIdentity()
		req.Resource = &i
	}
	if f.Action != nil {
		i := f.Action.GetIdentity()
		req.Action = &i
	}
	_, err := r.srv.RemoveSubjectPermission(ctx, req)
	return err
}
