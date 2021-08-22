package ent

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

type CreationAuditMixin struct {
	// We embed the `mixin.Schema` to avoid
	// implementing the rest of the methods.
	mixin.Schema
	mixin.Time
}

var _ ent.Mixin = (*CreationAuditMixin)(nil)

func (CreationAuditMixin) Fields() []ent.Field {
	return append(mixin.CreateTime{}.Fields(), field.String("created_by"))
}

type UpdateAuditMixin struct {
	// We embed the `mixin.Schema` to avoid
	// implementing the rest of the methods.
	mixin.Schema
}

var _ ent.Mixin = (*UpdateAuditMixin)(nil)

func (UpdateAuditMixin) Fields() []ent.Field {
	return append(mixin.UpdateTime{}.Fields(), field.String("created_by"))
}
