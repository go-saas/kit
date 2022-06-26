package biz

import (
	"github.com/google/uuid"
	gorm2 "github.com/go-saas/kit/pkg/gorm"
	"github.com/go-saas/saas/gorm"
)

type UserRole struct {
	gorm2.UIDBase
	gorm.MultiTenancy
	UserID uuid.UUID `gorm:"type:char(36)"`
	RoleID uuid.UUID `gorm:"type:char(36)"`
}
