package biz

import "github.com/google/uuid"

type UserLogin struct {
	UserId        uuid.UUID `gorm:"type:char(36);primaryKey" json:"user_id"`
	LoginProvider string    `gorm:"primaryKey" json:"login_provider"`
	ProviderKey   string    `gorm:"index" json:"provider_key"`
}
