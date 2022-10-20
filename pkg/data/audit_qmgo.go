package data

import (
	"context"
	"github.com/go-saas/kit/pkg/authn"
	"time"
)

func (model *AuditedModel) BeforeInsert(ctx context.Context) error {
	model.CreatedAt = time.Now().Local()
	if u, ok := authn.FromUserContext(ctx); ok {
		uid := u.GetId()
		model.CreatedBy = &uid
	}
	return nil
}
func (model *AuditedModel) BeforeUpdate(ctx context.Context) error {
	model.UpdatedAt = time.Now().Local()
	if u, ok := authn.FromUserContext(ctx); ok {
		uid := u.GetId()
		model.UpdatedBy = &uid
	}
	return nil
}
func (model *AuditedModel) BeforeUpsert(ctx context.Context) error {
	if model.CreatedAt.IsZero() {
		model.CreatedAt = time.Now().Local()
	}
	model.UpdatedAt = time.Now().Local()
	if u, ok := authn.FromUserContext(ctx); ok {
		uid := u.GetId()
		if model.CreatedBy == nil {
			model.CreatedBy = &uid
		}
		model.UpdatedBy = &uid
	}
	return nil
}
