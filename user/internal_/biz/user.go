package biz

import (
	"github.com/goxiaoy/go-saas-kit/pkg/gorm"
	gorm3 "github.com/goxiaoy/go-saas/gorm"
	concurrency "github.com/goxiaoy/gorm-concurrency"
	gorm2 "gorm.io/gorm"
	"time"
)

type User struct {
	gorm.UIDBase
	concurrency.Version `gorm:"type:char(36)"`
	gorm.AuditedModel
	gorm3.MultiTenancy

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
	Password         *string `json:"password"`
	TwoFactorEnabled bool    `json:"two_factor_enabled"`

	Roles []Role `gorm:"many2many:user_roles"`

	Location *string `json:"location"`
	Tags     *string `json:"tags"`

	// Avatar could be an id of asset or simple url
	Avatar   *string    `json:"avatar"`
	Birthday *time.Time `json:"birthday"`
	Gender   *string    `json:"gender" rql:"filter"`

	// SecondEmail back up email
	SecondEmail *string `json:"second_email"`

}
