package biz

import (
	"context"
	"github.com/goxiaoy/go-saas-kit/authorization/authorization"
	"github.com/goxiaoy/go-saas/seed"
)

type PermissionSeeder struct {
	permission authorization.PermissionManagementService
}

func NewPermissionSeeder(permission authorization.PermissionManagementService) *PermissionSeeder {
	return &PermissionSeeder{permission: permission}
}

func (p *PermissionSeeder) Seed(ctx context.Context, sCtx *seed.Context) error {
	return p.permission.AddGrant(ctx, authorization.NewEntityResource("*", "*"), authorization.ActionStr("*"), authorization.NewClientSubject("*"), authorization.EffectGrant)
}
