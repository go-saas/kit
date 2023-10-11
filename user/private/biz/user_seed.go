package biz

import (
	"context"
	"github.com/go-saas/saas/seed"
	"github.com/samber/lo"
)

const (
	Admin            = "admin"
	AdminUsernameKey = "admin_username"
	AdminEmailKey    = "admin_email"
	AdminPasswordKey = "admin_password"
	AdminUserId      = "admin_user_id"
)

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
	adminId := ""
	var admin *User
	var err error
	var ok bool
	var shouldCreate = false

	if adminId, ok = sCtx.Extra[AdminUserId].(string); ok {
		//attach existing user as tenant amin
		admin, err = u.um.FindByID(ctx, adminId)
	} else if adminUsername, ok = sCtx.Extra[AdminUsernameKey].(string); ok {
		shouldCreate = true
		admin, err = u.um.FindByName(ctx, adminUsername)
	} else if adminEmail, ok = sCtx.Extra[AdminEmailKey].(string); ok {
		shouldCreate = true
		admin, err = u.um.FindByEmail(ctx, adminEmail)
	}

	if err != nil {
		return err
	}
	adminPassword := ""
	adminPassword, _ = sCtx.Extra[AdminPasswordKey].(string)

	//seed admin
	if admin == nil && shouldCreate {
		//seed
		name := adminUsername
		admin = &User{
			Name: &name,
		}
		if len(adminUsername) > 0 {
			admin.Username = &adminUsername
		}
		if len(adminEmail) > 0 {
			admin.SetEmail(adminEmail, true)
		}
		if err = u.um.CreateWithPassword(ctx, admin, adminPassword, false); err != nil {
			return err
		}
	}
	if admin == nil {
		//can not create
		return nil
	}
	//add into role
	roles, err := u.um.GetRoles(ctx, admin.ID.String())
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
