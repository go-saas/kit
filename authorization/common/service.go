package common

import (
	"context"
	"github.com/goxiaoy/go-saas-kit/auth/current"
)

type AuthorizationService interface {
	Check(ctx context.Context, resource Resource, action Action, namespace Namespace, subject Subject) (AuthorizationResult, error)
	CheckCurrentUser(ctx context.Context, resource Resource, action Action, namespace Namespace) (AuthorizationResult, error)
}

type AuthorizationResult struct {
	Allowed      bool
	Requirements []Requirement
}

func NewAllowAuthorizationResult() AuthorizationResult {
	return AuthorizationResult{Allowed: true}
}

func NewDisAllowAuthorizationResult(requirements []Requirement) AuthorizationResult {
	return AuthorizationResult{Allowed: false, Requirements: requirements}
}

type AuthenticationAuthorizationService struct {
}

var _ AuthorizationService = (*AuthenticationAuthorizationService)(nil)

func NewAuthenticationAuthorizationService() *AuthenticationAuthorizationService {
	return &AuthenticationAuthorizationService{}
}

func (a *AuthenticationAuthorizationService) Check(ctx context.Context, resource Resource, action Action, namespace Namespace, subject Subject) (AuthorizationResult, error) {
	if always, ok := FromAlwaysAuthorizationContext(ctx); ok {
		if always {
			return NewAllowAuthorizationResult(), nil
		} else {
			return NewDisAllowAuthorizationResult(nil), nil
		}
	}
	var userId string
	if us, ok := subject.(*UserSubject); ok {
		userId = us.GetIdentity()
	}
	if userId == "" {
		return NewDisAllowAuthorizationResult([]Requirement{AuthenticationRequirement}), nil
	} else {
		return NewAllowAuthorizationResult(), nil
	}
}

func (a *AuthenticationAuthorizationService) CheckCurrentUser(ctx context.Context, resource Resource, action Action, namespace Namespace) (AuthorizationResult, error) {
	var userId string
	if userInfo, ok := current.FromUserContext(ctx); ok {
		userId = userInfo.GetId()
	}
	return a.Check(ctx, namespace, resource, action, NewUserSubject(userId))
}
