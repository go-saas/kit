package biz

import (
	"github.com/google/uuid"
	"github.com/goxiaoy/go-saas/gorm"
	gorm2 "github.com/goxiaoy/go-saas-kit/pkg/gorm"
)

type UserRole struct {
	gorm2.UIDBase
	gorm.MultiTenancy
	UserID uuid.UUID `gorm:"type:char(36)"`
	RoleID uuid.UUID `gorm:"type:char(36)"`
}
