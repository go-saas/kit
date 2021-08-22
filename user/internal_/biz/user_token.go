package biz

import "github.com/google/uuid"

type UserToken struct {
	UserId        uuid.UUID `gorm:"type:char(36);primaryKey" json:"user_id"`
	LoginProvider string    `gorm:"primaryKey" json:"login_provider"`
	Name          string    `gorm:"primaryKey" json:"name"`
	Value         string    `json:"value"`
}
