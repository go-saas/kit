package remote

import (
	"context"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authz"
	v1 "github.com/goxiaoy/go-saas-kit/user/api/permission/v1"
	"github.com/goxiaoy/go-saas/common"
)

type PermissionChecker struct {
	client v1.PermissionServiceClient
}

func NewRemotePermissionChecker(client v1.PermissionServiceClient) authz.PermissionChecker {
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

func (r *PermissionChecker) IsGrant(ctx context.Context, resource authz.Resource, action authz.Action, subjects ...authz.Subject) (authz.Effect, error) {
	tenantInfo := common.FromCurrentTenant(ctx)
	return r.IsGrantTenant(ctx, resource, action, tenantInfo.GetId(), subjects...)
}
