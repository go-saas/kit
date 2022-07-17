package biz

import (
	"context"
	"github.com/go-saas/kit/event"
	v12 "github.com/go-saas/kit/user/event/v1"
	"github.com/go-saas/saas/gorm"
	"github.com/google/uuid"
)

type UserRole struct {
	gorm.MultiTenancy `gorm:"primaryKey"`
	UserID            uuid.UUID `gorm:"type:char(36);primaryKey"`
	RoleID            uuid.UUID `gorm:"type:char(36);primaryKey"`
}

func NewUserRoleChangeEventHandler(um *UserManager) event.ConsumerHandler {
	msg := &v12.UserRoleChangeEvent{}
	return event.ProtoHandler[*v12.UserRoleChangeEvent](msg, event.HandlerFuncOf[*v12.UserRoleChangeEvent](func(ctx context.Context, msg *v12.UserRoleChangeEvent) error {
		return um.RemoveUserRoleCache(ctx, msg.UserId)
	}))
}
