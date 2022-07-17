package biz

import (
	"context"
	"errors"
	"github.com/go-saas/kit/pkg/authz/authz"
	"github.com/go-saas/saas"

	"github.com/go-saas/saas/seed"
)

type PermissionSeeder struct {
	permission authz.PermissionManagementService
	checker    authz.PermissionChecker
	rm         *RoleManager
}

func NewPermissionSeeder(permission authz.PermissionManagementService, checker authz.PermissionChecker, rm *RoleManager) *PermissionSeeder {
	return &PermissionSeeder{permission: permission, checker: checker, rm: rm}
}

func (p *PermissionSeeder) Seed(ctx context.Context, sCtx *seed.Context) error {

	tenantInfo, _ := saas.FromCurrentTenant(ctx)
	err := authz.EnsureGrant(ctx, p.permission, p.checker,
		authz.NewEntityResource("*", "*"),
		authz.ActionStr("*"),
		authz.NewClientSubject("*"),
		"*")
	if err != nil {
		return err
	}
	//find admin role
	adminRole, err := p.rm.FindByName(ctx, Admin)
	if err != nil {
		return err
	}
	if adminRole == nil {
		return errors.New("admin role not found")
	}
	permissionTenant := tenantInfo.GetId()
	if len(permissionTenant) == 0 {
		//host
		//TODO this make admin of host to be a super admin across all tenants
		permissionTenant = "*"
	}
	if err := authz.EnsureGrant(ctx, p.permission, p.checker,
		authz.NewEntityResource("*", "*"),
		authz.ActionStr("*"),
		authz.NewRoleSubject(adminRole.ID.String()),
		permissionTenant); err != nil {
		return err
	}

	return nil
}
