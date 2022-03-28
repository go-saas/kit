package authz

import (
	"context"
)

type PermissionChecker interface {
	IsGrantTenant(ctx context.Context, requirement RequirementList, tenantID string, subjects ...Subject) ([]Effect, error)
}
