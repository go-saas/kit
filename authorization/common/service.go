package common

import "context"

type AuthorizationService interface {
	Check(ctx context.Context, namespace Namespace, resource Resource, action Action, subject Subject) (AuthorizationResult, error)
}

type AuthorizationResult interface {
	Allowed() bool
	GetRequirements() []Requirement
}
