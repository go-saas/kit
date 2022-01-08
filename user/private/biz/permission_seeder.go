package biz

import (
	"context"
	authorization2 "github.com/goxiaoy/go-saas-kit/pkg/authz/authorization"
	"github.com/goxiaoy/go-saas/seed"
)

type PermissionSeeder struct {
	permission authorization2.PermissionManagementService
}

func NewPermissionSeeder(permission authorization2.PermissionManagementService) *PermissionSeeder {
	return &PermissionSeeder{permission: permission}
}

func (p *PermissionSeeder) Seed(ctx context.Context, sCtx *seed.Context) error {
	return p.permission.AddGrant(ctx, authorization2.NewEntityResource("*", "*"), authorization2.ActionStr("*"), authorization2.NewClientSubject("*"), authorization2.EffectGrant)
}
