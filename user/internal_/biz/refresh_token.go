package biz

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type RefreshToken struct {
	Token     string    `gorm:"primaryKey"`
	UserId    uuid.UUID `gorm:"type:char(36);primaryKey" json:"user_id"`
	Expires   time.Time
	Ip        string
	UserAgent string
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
