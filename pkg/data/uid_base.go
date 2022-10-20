package data

import "github.com/google/uuid"

type UIDBase struct {
	ID uuid.UUID `gorm:"type:char(36)" json:"id"`
}
