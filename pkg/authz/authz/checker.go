package authz

import (
	"context"
)

type PermissionChecker interface {
	IsGrantTenant(ctx context.Context, resource Resource, action Action, tenantID string, subjects ...Subject) (Effect, error)
}
