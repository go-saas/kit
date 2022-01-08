package biz

import (
	"context"
	authorization2 "github.com/goxiaoy/go-saas-kit/pkg/authz/authorization"
	v1 "github.com/goxiaoy/go-saas-kit/user/api/role/v1"
	"github.com/goxiaoy/go-saas/seed"
)

type PermissionSeeder struct {
	permission  authorization2.PermissionManagementService
	roleService v1.RoleServiceClient
}

func NewPermissionSeeder(permission authorization2.PermissionManagementService, roleService v1.RoleServiceClient) *PermissionSeeder {
	return &PermissionSeeder{permission: permission, roleService: roleService}
}

func (p *PermissionSeeder) Seed(ctx context.Context, sCtx *seed.Context) error {
	admin, err := p.roleService.GetRole(ctx, &v1.GetRoleRequest{Name: "admin"})
	if err != nil {
		return err
	}
	if err := p.permission.AddGrant(ctx, authorization2.NewEntityResource("*", "*"),
		authorization2.ActionStr("*"), authorization2.NewClientSubject("*"), authorization2.EffectGrant); err != nil {
		return err
	}

	if err := p.permission.AddGrant(ctx, authorization2.NewEntityResource("*", "*"),
		authorization2.ActionStr("*"), authorization2.NewRoleSubject(admin.Id), authorization2.EffectGrant); err != nil {
		return err
	}

	return nil
}
