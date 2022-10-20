package data

import (
	"context"
	"github.com/google/uuid"
)

func (u *UIDBase) BeforeInsert(ctx context.Context) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
func (u *UIDBase) BeforeUpsert(ctx context.Context) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
