package keto

import (
	"context"
	"github.com/goxiaoy/go-saas-kit/authorization/common"
	"github.com/ory/keto/proto/ory/keto/acl/v1alpha1"
)

type KetoAuthorizationService struct {
	*common.AuthenticationAuthorizationService
	client acl.CheckServiceClient
}

var _ common.AuthorizationService = (*KetoAuthorizationService)(nil)

func NewKetoAuthorizationService(client acl.CheckServiceClient) *KetoAuthorizationService {
	return &KetoAuthorizationService{common.NewAuthenticationAuthorizationService(),client}
}

func (k KetoAuthorizationService) Check(ctx context.Context, resource common.Resource, action common.Action, namespace common.Namespace, subject common.Subject) (common.AuthorizationResult, error) {
	req:=&acl.CheckRequest{}
	if namespace!=nil{
		req.Namespace=namespace.GetIdentity()
	}
	if resource!=nil{
		req.Object=resource.GetIdentity()
	}
	if action!=nil{
		req.Relation=action.GetIdentity()
	}
	if subject!=nil{
		req.Subject=&acl.Subject{Ref:&acl.Subject_Id{Id: subject.GetIdentity()}}
	}
	//TODO get snaptoken from context
	resp,err:=k.client.Check(ctx,req)
	if err!=nil{
		return common.NewDisAllowAuthorizationResult(nil),err
	}
	if resp.Allowed{
		return common.NewAllowAuthorizationResult(),nil
	}else{
		return common.NewDisAllowAuthorizationResult([]common.Requirement{
			//TODO requirement
		}),nil
	}
}

