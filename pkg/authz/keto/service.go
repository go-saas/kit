package keto

import (
	"context"
	authorization2 "github.com/goxiaoy/go-saas-kit/pkg/authz/authorization"
	"github.com/ory/keto/proto/ory/keto/acl/v1alpha1"
)

type PermissionChecker struct {
	client acl.CheckServiceClient
}

var _ authorization2.PermissionChecker = (*PermissionChecker)(nil)

func NewPermissionChecker(client acl.CheckServiceClient) *PermissionChecker {
	return &PermissionChecker{client}
}

func (k *PermissionChecker) IsGrant(ctx context.Context, resource authorization2.Resource, action authorization2.Action, subjects ...authorization2.Subject) (authorization2.Effect, error) {
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
			return authorization2.EffectUnknown, err
		}
		//TODO keto do not support multiple subjects
		if !resp.Allowed {
			return authorization2.EffectForbidden, nil
		}

	}
	return authorization2.EffectGrant, nil
}
