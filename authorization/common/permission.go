package common

import (
	"context"
	"fmt"
	"sync"
)

type PermissionManagementService interface {
	AddGrant(ctx context.Context, resource Resource, action Action, subject Subject, grantType GrantType) error
}

type PermissionChecker interface {
	IsGrant(ctx context.Context, resource Resource, action Action, subject Subject) (GrantType, error)
}

type permissionBean struct {
	namespace string
	resource  string
	action    string
	subject   string
	grantType GrantType
}

func newPermissionBean(resource Resource, action Action, subject Subject, grantType GrantType) permissionBean {
	return permissionBean{
		namespace: resource.GetNamespace(),
		resource:  resource.GetIdentity(),
		action:    action.GetIdentity(),
		subject:   fmt.Sprintf("%s/%s", subject.GetName(), subject.GetIdentity()),
		grantType: grantType,
	}
}

type PermissionService struct {
	v   []permissionBean
	mux sync.Mutex
}

var _ PermissionManagementService = (*PermissionService)(nil)
var _ PermissionChecker = (*PermissionService)(nil)

func NewPermissionService() *PermissionService {
	return &PermissionService{}
}

func (p *PermissionService) IsGrant(ctx context.Context, resource Resource, action Action, subject Subject) (GrantType, error) {
	p.mux.Lock()
	defer p.mux.Unlock()
	var anyAllow bool
	for _, bean := range p.v {
		//TODO regex match
		if bean.namespace == resource.GetNamespace() && bean.resource == resource.GetIdentity() && bean.action == action.GetIdentity() && bean.subject == subject.GetIdentity() {
			if bean.grantType == GrantTypeDisallow {
				return GrantTypeDisallow, nil
			}
			if bean.grantType == GrantTypeAllow {
				anyAllow = true
			}
		}
	}
	if anyAllow {
		return GrantTypeAllow, nil
	}
	return GrantTypeUnknown, nil
}

func (p *PermissionService) AddGrant(ctx context.Context, resource Resource, action Action, subject Subject, grantType GrantType) error {
	p.mux.Lock()
	defer p.mux.Unlock()
	p.v = append(p.v, newPermissionBean(resource, action, subject, grantType))
	return nil
}
