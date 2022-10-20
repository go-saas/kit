package data

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (u *UIDBase) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
