package biz

import (
	"context"
	"github.com/go-saas/kit/pkg/data"
	kitgorm "github.com/go-saas/kit/pkg/gorm"
	v1 "github.com/go-saas/kit/user/api/user/v1"
	"github.com/go-saas/saas"
	concurrency "github.com/goxiaoy/gorm-concurrency/v2"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"time"
)

type User struct {
	kitgorm.UIDBase        `json:",squash"`
	concurrency.HasVersion `gorm:"type:char(36)"`
	kitgorm.AuditedModel
	kitgorm.AggRoot

	DeletedAt gorm.DeletedAt `gorm:"index"`

	Name      *string `json:"name"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`

	Username *string `json:"username"`
	// NormalizedUsername uppercase normalized userName
	NormalizedUsername *string `json:"normalized_username" gorm:"index:,size:200"`

	// Phone
	Phone          *string `json:"phone" gorm:"index:,size:200"`
	PhoneConfirmed bool    `json:"phone_confirmed"`

	// Email
	Email *string `json:"email"`
	// NormalizedEmail uppercase normalized email
	NormalizedEmail *string `json:"normalized_email" gorm:"index:,size:200"`
	EmailConfirmed  bool    `json:"email_confirmed"`

	// Password hashed
	Password *string `json:"password"`

	//Security
	AccessFailedCount int        `json:"accessFailedCount"`
	LastLoginAttempt  *time.Time `json:"lastLoginAttempt"`
	LockoutEndDateUtc *time.Time `json:"lockoutEndDateUtc"`

	//2FA
	TwoFactorEnabled bool `json:"two_factor_enabled"`

	Roles []Role `gorm:"many2many:user_roles"`

	Location *string `json:"location"`
	Tags     *string `json:"tags"`

	// Avatar could be an id of asset or simple url
	Avatar   *string    `json:"avatar"`
	Birthday *time.Time `json:"birthday"`
	Gender   *string    `json:"gender"`

	Tenants []UserTenant `json:"tenants"`

	Extra data.JSONMap
	//creation tenant
	CreatedTenant *string `json:"created_tenant"`
}

type UserRepo interface {
	data.Repo[User, string, *v1.ListUsersRequest]

	ListAdmin(ctx context.Context, query *v1.AdminListUsersRequest) ([]*User, error)
	CountAdmin(ctx context.Context, query *v1.AdminListUsersRequest) (total int64, filtered int64, err error)

	FindByID(ctx context.Context, id string) (*User, error)
	FindByName(ctx context.Context, name string) (*User, error)
	FindByPhone(ctx context.Context, phone string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)

	AddLogin(ctx context.Context, user *User, userLogin *UserLogin) error
	RemoveLogin(ctx context.Context, user *User, loginProvider string, providerKey string) error
	ListLogin(ctx context.Context, user *User) ([]*UserLogin, error)
	FindByLogin(ctx context.Context, loginProvider string, providerKey string) (*User, error)

	SetToken(ctx context.Context, user *User, loginProvider string, name string, value string) error
	RemoveToken(ctx context.Context, user *User, loginProvider string, name string) error
	GetToken(ctx context.Context, user *User, loginProvider string, name string) (*string, error)

	GetRoles(ctx context.Context, userId string) ([]Role, error)
	UpdateRoles(ctx context.Context, user *User, roles []Role) error
	AddToRole(ctx context.Context, user *User, role *Role) error
	RemoveFromRole(ctx context.Context, user *User, role *Role) error
}

func (u *User) SetPhone(phone string, confirmed bool) {
	u.Phone = &phone
	u.PhoneConfirmed = confirmed
}

func (u *User) SetEmail(email string, confirmed bool) {
	u.Email = &email
	u.EmailConfirmed = confirmed
}

func (u *User) CheckInCurrentTenant(ctx context.Context) error {
	ct, _ := saas.FromCurrentTenant(ctx)
	_, isIn := lo.Find(u.Tenants, func(tenant UserTenant) bool {
		return tenant.TenantId == ct.GetId()
	})
	if !isIn {
		return v1.ErrorUserNotFoundLocalized(ctx, nil, nil)
	}
	return nil
}
