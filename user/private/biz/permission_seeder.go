package biz

import (
	"context"
	"errors"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authorization"
	"github.com/goxiaoy/go-saas/common"
	"github.com/goxiaoy/go-saas/seed"
)

type PermissionSeeder struct {
	permission authorization.PermissionManagementService
	checker    authorization.PermissionChecker
	rm         *RoleManager
}

func NewPermissionSeeder(permission authorization.PermissionManagementService, checker authorization.PermissionChecker, rm *RoleManager) *PermissionSeeder {
	return &PermissionSeeder{permission: permission, checker: checker, rm: rm}
}

func (p *PermissionSeeder) Seed(ctx context.Context, sCtx *seed.Context) error {

	tenantInfo := common.FromCurrentTenant(ctx)
	err := authorization.EnsureGrant(ctx, p.permission, p.checker,
		authorization.NewEntityResource("*", "*"),
		authorization.ActionStr("*"),
		authorization.NewClientSubject("*"),
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
	err = authorization.EnsureGrant(ctx, p.permission, p.checker,
		authorization.NewEntityResource("*", "*"),
		authorization.ActionStr("*"),
		authorization.NewRoleSubject(adminRole.ID.String()),
		tenantInfo.GetId())
	if err != nil {
		return err
	}

	return nil
}
