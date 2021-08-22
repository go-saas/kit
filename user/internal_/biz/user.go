package biz

import (
	gorm3 "github.com/goxiaoy/go-saas/gorm"
	concurrency "github.com/goxiaoy/gorm-concurrency"
	gorm2 "gorm.io/gorm"
	"github.com/goxiaoy/go-saas-kit/pkg/gorm"
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

	UserName *string `json:"user_name" gorm:"index"`
	// NormalizedUserName uppercase normalized userName
	NormalizedUserName *string `json:"normalized_user_name" gorm:"index"`

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

	// Avatar could be a id of asset or simple url
	Avatar   *string    `json:"avatar"`
	Birthday *time.Time `json:"birthday"`
	Gender   *string    `json:"gender" rql:"filter"`

	// SecondEmail back up email
	SecondEmail *string `json:"second_email"`

}
