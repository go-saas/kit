package biz

import (
	"context"
	"github.com/goxiaoy/go-saas-kit/authorization/authorization"
	v1 "github.com/goxiaoy/go-saas-kit/user/api/role/v1"
	"github.com/goxiaoy/go-saas/seed"
)

type PermissionSeeder struct {
	permission  authorization.PermissionManagementService
	roleService v1.RoleServiceClient
}

func NewPermissionSeeder(permission authorization.PermissionManagementService, roleService v1.RoleServiceClient) *PermissionSeeder {
	return &PermissionSeeder{permission: permission, roleService: roleService}
}

func (p *PermissionSeeder) Seed(ctx context.Context, sCtx *seed.Context) error {
	admin, err := p.roleService.GetRole(ctx, &v1.GetRoleRequest{Name: "admin"})
	if err != nil {
		return err
	}
	return p.permission.AddGrant(ctx, authorization.NewEntityResource("*", admin.Id), authorization.ActionStr("*"), authorization.NewClientSubject("*"), authorization.EffectGrant)
}
