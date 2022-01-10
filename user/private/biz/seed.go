package biz

import (
	"context"
	"github.com/ahmetb/go-linq/v3"
	authorization2 "github.com/goxiaoy/go-saas-kit/pkg/authz/authorization"
	"github.com/goxiaoy/go-saas/seed"
)

const Admin = "admin"
const AdminUsernameKey = "admin_username"
const AdminPasswordKey = "admin_password"

type RoleSeed struct {
	rm         *RoleManager
	permission authorization2.PermissionManagementService
}

func NewRoleSeed(roleMgr *RoleManager, permission authorization2.PermissionManagementService) *RoleSeed {
	return &RoleSeed{rm: roleMgr, permission: permission}
}

func (r *RoleSeed) Seed(ctx context.Context, sCtx *seed.Context) error {
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
		if role.Name == Admin {
			r.permission.AddGrant(ctx, authorization2.NewEntityResource("*", "*"), authorization2.ActionStr("*"), authorization2.NewRoleSubject(role.ID.String()), authorization2.EffectGrant)
		}
	}
	return nil
}

type UserSeed struct {
	um *UserManager
	rm *RoleManager
}

func NewUserSeed(um *UserManager, rm *RoleManager) *UserSeed {
	return &UserSeed{um: um, rm: rm}
}
func (u *UserSeed) Seed(ctx context.Context, sCtx *seed.Context) error {
	adminUsername := ""
	adminUsername, _ = sCtx.Extra[AdminUsernameKey].(string)

	adminPassword := ""
	adminPassword, _ = sCtx.Extra[AdminPasswordKey].(string)

	admin, err := u.um.FindByName(ctx, adminUsername)
	if err != nil {
		return err
	}
	if admin == nil {
		//seed
		name := adminUsername
		admin = &User{
			Name:     &name,
			Username: &name,
		}
		if err = u.um.CreateWithPassword(ctx, admin, adminPassword); err != nil {
			return err
		}
	}
	//add into role
	roles, err := u.um.GetRoles(ctx, admin)
	if err != nil {
		return err
	}
	if find := linq.From(roles).AnyWithT(func(r *Role) bool {
		return r.Name == Admin
	}); !find {
		adminRole, err := u.rm.FindByName(ctx, Admin)
		if err != nil {
			return err
		}
		if err := u.um.AddToRole(ctx, admin, adminRole); err != nil {
			return err
		}
	}
	return nil
}