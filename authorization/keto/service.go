package keto

import (
	"context"
	"fmt"
	"github.com/goxiaoy/go-saas-kit/authorization/common"
	"github.com/ory/keto/proto/ory/keto/acl/v1alpha1"
)

type PermissionChecker struct {
	client acl.CheckServiceClient
}

var _ common.PermissionChecker = (*PermissionChecker)(nil)

func NewPermissionChecker(client acl.CheckServiceClient) *PermissionChecker {
	return &PermissionChecker{client}
}

func (k *PermissionChecker) IsGrant(ctx context.Context, resource common.Resource, action common.Action, subject common.Subject) (common.GrantType, error) {
	req := &acl.CheckRequest{}

	req.Namespace = resource.GetNamespace()
	req.Object = resource.GetIdentity()

	if action != nil {
		req.Relation = action.GetIdentity()
	}
	if subject != nil {
		req.Subject = &acl.Subject{Ref: &acl.Subject_Id{Id: fmt.Sprintf("%s/%s", subject.GetName(), subject.GetIdentity())}}
	}
	//TODO get snaptoken from context
	resp, err := k.client.Check(ctx, req)
	if err != nil {
		return common.GrantTypeUnknown, err
	}
	if resp.Allowed {
		return common.GrantTypeAllow, nil
	} else {
		return common.GrantTypeUnknown, nil
	}
}
