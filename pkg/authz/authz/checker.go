package authz

import (
	"context"
)

type PermissionChecker interface {
	IsGrant(ctx context.Context, resource Resource, action Action, subjects ...Subject) (Effect, error)
	IsGrantTenant(ctx context.Context, resource Resource, action Action, tenantID string, subjects ...Subject) (Effect, error)
}
