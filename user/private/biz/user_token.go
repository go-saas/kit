package biz

import (
	"context"
	"github.com/go-saas/kit/pkg/gorm"
	"github.com/google/uuid"
	gorm2 "gorm.io/gorm"
)

const (
	InternalLoginProvider     string = "internal"
	InternalRememberTokenName string = "remember"
)

//UserToken stores external login token
type UserToken struct {
	gorm.AuditedModel
	DeletedAt     gorm2.DeletedAt `gorm:"index"`
	UserId        uuid.UUID       `gorm:"type:char(36);primaryKey" json:"user_id"`
	LoginProvider string          `gorm:"primaryKey" json:"login_provider"`
	Name          string          `gorm:"primaryKey" json:"name"`
	Value         string          `json:"value"`
}

type UserTokenRepo interface {
	FindByUserIdAndLoginProvider(ctx context.Context, userId, loginProvider string) ([]*UserToken, error)
	FindByUserIdAndLoginProviderAndName(ctx context.Context, userId, loginProvider, name string) (*UserToken, error)
	DeleteByUserIdAndLoginProvider(ctx context.Context, userId, loginProvider string) error
	DeleteByUserIdAndLoginProviderAndName(ctx context.Context, userId, loginProvider, name string) error
	Create(ctx context.Context, userId, loginProvider, name, value string) (*UserToken, error)
}
