package api

import (
	"context"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authorization"
	v1 "github.com/goxiaoy/go-saas-kit/user/api/permission/v1"
	"github.com/goxiaoy/go-saas/common"
)

type RemotePermissionChecker struct {
	client v1.PermissionServiceClient
}

func NewRemotePermissionChecker(client v1.PermissionServiceClient) authorization.PermissionChecker {
	return &RemotePermissionChecker{
		client: client,
	}
}

func (r *RemotePermissionChecker) IsGrantTenant(ctx context.Context, resource authorization.Resource, action authorization.Action, tenantID string, subjects ...authorization.Subject) (authorization.Effect, error) {
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
		return authorization.EffectForbidden, err
	}
	if grant.Effect == v1.Effect_GRANT {
		return authorization.EffectGrant, nil
	}
	return authorization.EffectForbidden, nil
}

func (r *RemotePermissionChecker) IsGrant(ctx context.Context, resource authorization.Resource, action authorization.Action, subjects ...authorization.Subject) (authorization.Effect, error) {
	tenantInfo := common.FromCurrentTenant(ctx)
	return r.IsGrantTenant(ctx, resource, action, tenantInfo.GetId(), subjects...)
}
