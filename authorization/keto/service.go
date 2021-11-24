package keto

import (
	"context"
	"github.com/goxiaoy/go-saas-kit/authorization/common"
	"github.com/ory/keto/proto/ory/keto/acl/v1alpha1"
)

type AuthorizationService struct {
	*common.AuthenticationAuthorizationService
	client acl.CheckServiceClient
}

var _ common.AuthorizationService = (*AuthorizationService)(nil)

func NewAuthorizationService(client acl.CheckServiceClient) *AuthorizationService {
	return &AuthorizationService{common.NewAuthenticationAuthorizationService(), client}
}

func (k AuthorizationService) Check(ctx context.Context, resource common.Resource, action common.Action, subject common.Subject) (common.AuthorizationResult, error) {
	req := &acl.CheckRequest{}

	req.Namespace = resource.GetNamespace()
	req.Object = resource.GetIdentity()

	if action != nil {
		req.Relation = action.GetIdentity()
	}
	if subject != nil {
		req.Subject = &acl.Subject{Ref: &acl.Subject_Id{Id: subject.GetIdentity()}}
	}
	//TODO get snaptoken from context
	resp, err := k.client.Check(ctx, req)
	if err != nil {
		return common.NewDisallowAuthorizationResult(nil), err
	}
	if resp.Allowed {
		return common.NewAllowAuthorizationResult(), nil
	} else {
		return common.NewDisallowAuthorizationResult([]common.Requirement{
			common.NewRequirement(resource, action, subject, ""),
		}), nil
	}
}
