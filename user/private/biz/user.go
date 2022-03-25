package biz

import (
	"context"
	"github.com/goxiaoy/go-saas-kit/pkg/data"
	"github.com/goxiaoy/go-saas-kit/pkg/gorm"
	v1 "github.com/goxiaoy/go-saas-kit/user/api/user/v1"
	concurrency "github.com/goxiaoy/gorm-concurrency"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	gorm2 "gorm.io/gorm"
	"time"
)

type User struct {
	gorm.UIDBase        `json:",squash"`
	concurrency.Version `gorm:"type:char(36)"`
	gorm.AuditedModel

	DeletedAt gorm2.DeletedAt `gorm:"index"`

	Name      *string `json:"name"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`

	Username *string `json:"username" gorm:"index"`
	// NormalizedUsername uppercase normalized userName
	NormalizedUsername *string `json:"normalized_username" gorm:"index"`

	// Phone
	Phone          *string `json:"phone" gorm:"index"`
	PhoneConfirmed bool    `json:"phone_confirmed"`

	// Email
	Email *string `json:"email" gorm:"index"`
	// NormalizedEmail uppercase normalized email
	NormalizedEmail *string `json:"normalized_email" gorm:"index"`
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
	List(ctx context.Context, query *v1.ListUsersRequest) ([]*User, error)
	Count(ctx context.Context, query *v1.UserFilter) (total int64, filtered int64, err error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User, p *fieldmaskpb.FieldMask) error
	Delete(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id string) (*User, error)
	FindByName(ctx context.Context, name string) (*User, error)
	FindByPhone(ctx context.Context, phone string) (*User, error)

	AddLogin(ctx context.Context, user *User, userLogin *UserLogin) error
	RemoveLogin(ctx context.Context, user *User, loginProvider string, providerKey string) error
	ListLogin(ctx context.Context, user *User) ([]*UserLogin, error)

	FindByLogin(ctx context.Context, loginProvider string, providerKey string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	SetToken(ctx context.Context, user *User, loginProvider string, name string, value string) error
	RemoveToken(ctx context.Context, user *User, loginProvider string, name string) error
	GetToken(ctx context.Context, user *User, loginProvider string, name string) (*string, error)
	GetRoles(ctx context.Context, user *User) ([]*Role, error)
	UpdateRoles(ctx context.Context, user *User, roles []*Role) error
	AddToRole(ctx context.Context, user *User, role *Role) error
	RemoveFromRole(ctx context.Context, user *User, role *Role) error
}
