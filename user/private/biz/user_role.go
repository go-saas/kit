package biz

import (
	"context"
	"github.com/go-saas/kit/event"
	gorm2 "github.com/go-saas/kit/pkg/gorm"
	v12 "github.com/go-saas/kit/user/event/v1"
	"github.com/go-saas/saas/gorm"
	"github.com/google/uuid"
)

type UserRole struct {
	gorm2.UIDBase
	gorm.MultiTenancy
	UserID uuid.UUID `gorm:"type:char(36)"`
	RoleID uuid.UUID `gorm:"type:char(36)"`
}

type UserRoleChangeEventHandler event.ConsumerHandler

func NewUserRoleChangeEventHandler(um *UserManager) UserRoleChangeEventHandler {
	msg := &v12.UserRoleChangeEvent{}
	return event.ProtoHandler[*v12.UserRoleChangeEvent](msg, event.HandlerFuncOf[*v12.UserRoleChangeEvent](func(ctx context.Context, msg *v12.UserRoleChangeEvent) error {
		return um.RemoveUserRoleCache(ctx, msg.UserId)
	}))
}
