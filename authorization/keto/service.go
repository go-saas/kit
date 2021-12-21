package keto

import (
	"context"
	"github.com/goxiaoy/go-saas-kit/authorization/authorization"
	"github.com/ory/keto/proto/ory/keto/acl/v1alpha1"
)

type PermissionChecker struct {
	client acl.CheckServiceClient
}

var _ authorization.PermissionChecker = (*PermissionChecker)(nil)

func NewPermissionChecker(client acl.CheckServiceClient) *PermissionChecker {
	return &PermissionChecker{client}
}

func (k *PermissionChecker) IsGrant(ctx context.Context, resource authorization.Resource, action authorization.Action, subjects ...authorization.Subject) (authorization.Effect, error) {
	req := &acl.CheckRequest{}

	req.Namespace = resource.GetNamespace()
	req.Object = resource.GetIdentity()

	if action != nil {
		req.Relation = action.GetIdentity()
	}
	for _, subject := range subjects {
		req.Subject = &acl.Subject{Ref: &acl.Subject_Id{Id: subject.GetIdentity()}}
		//TODO get snaptoken from context
		resp, err := k.client.Check(ctx, req)
		if err != nil {
			return authorization.EffectUnknown, err
		}
		//TODO keto do not support multiple subjects
		if !resp.Allowed {
			return authorization.EffectForbidden, nil
		}

	}
	return authorization.EffectGrant, nil
}
