package gorm

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UIDBase struct {
	ID uuid.UUID `gorm:"type:char(36)" json:"id"`
}

func (u *UIDBase) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
