package biz

import (
	"context"
	"github.com/go-saas/kit/pkg/authz/authz"
	"github.com/go-saas/saas/seed"
)

type RoleSeed struct {
	rm         *RoleManager
	permission authz.PermissionManagementService
}

func NewRoleSeed(roleMgr *RoleManager, permission authz.PermissionManagementService) *RoleSeed {
	return &RoleSeed{rm: roleMgr, permission: permission}
}

func (r *RoleSeed) Seed(ctx context.Context, _ *seed.Context) error {
	seedRoles := []*Role{
		{
			Name:        Admin,
			IsPreserved: true,
		},
	}
	for _, sr := range seedRoles {
		role, err := r.rm.FindByName(ctx, sr.Name)
		if err != nil {
			return err
		}
		if role == nil {
			if err := r.rm.Create(ctx, sr); err != nil {
				return err
			}
		}
	}
	return nil
}
