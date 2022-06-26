package biz

import (
	"context"
	"github.com/go-saas/kit/pkg/authz/authz"
	"github.com/go-saas/saas/seed"
	"github.com/samber/lo"
)

const (
	Admin            = "admin"
	AdminUsernameKey = "admin_username"
	AdminEmailKey    = "admin_email"
	AdminPasswordKey = "admin_password"
)

type RoleSeed struct {
	rm         *RoleManager
	permission authz.PermissionManagementService
}

func NewRoleSeed(roleMgr *RoleManager, permission authz.PermissionManagementService) *RoleSeed {
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
	adminEmail := ""
	var admin *User
	var err error
	var ok bool
	if adminUsername, ok = sCtx.Extra[AdminUsernameKey].(string); ok {
		admin, err = u.um.FindByName(ctx, adminUsername)
	} else if adminEmail, ok = sCtx.Extra[AdminEmailKey].(string); ok {
		admin, err = u.um.FindByEmail(ctx, adminEmail)
	}
	if err != nil {
		return err
	}
	adminPassword := ""
	adminPassword, _ = sCtx.Extra[AdminPasswordKey].(string)

	if len(adminUsername) == 0 && len(adminEmail) == 0 {
		//can not seed admin
		return nil
	}
	//seed admin
	if admin == nil {
		//seed
		name := adminUsername
		admin = &User{
			Name: &name,
		}
		if len(adminUsername) > 0 {
			admin.Username = &adminUsername
		}
		if len(adminEmail) > 0 {
			admin.Email = &adminEmail
		}
		if err = u.um.CreateWithPassword(ctx, admin, adminPassword, false); err != nil {
			return err
		}
	}
	//add into role
	roles, err := u.um.GetRoles(ctx, admin)
	if err != nil {
		return err
	}
	if find := lo.ContainsBy(roles, func(r Role) bool {
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

	//add into tenant
	if err := u.um.JoinTenant(ctx, admin.UIDBase.ID.String(), sCtx.TenantId); err != nil {
		return err
	}

	return nil
}
