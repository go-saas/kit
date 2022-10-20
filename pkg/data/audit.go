package data

import (
	"fmt"
	"time"
)

// AuditedModel make Model Auditable, embed `audited.AuditedModel` into your model as anonymous field to make the model auditable
//
//	type User struct {
//	  data.AuditedModel
//	}
type AuditedModel struct {
	CreatedBy *string
	UpdatedBy *string
	CreatedAt time.Time `gorm:"timestamp"`
	UpdatedAt time.Time `gorm:"timestamp"`
}

var _ auditableInterface = (*AuditedModel)(nil)

// SetCreatedBy set created by
func (model *AuditedModel) SetCreatedBy(createdBy interface{}) {
	if createdBy == nil {
		model.CreatedBy = nil
	} else {
		v := fmt.Sprintf("%v", createdBy)
		model.CreatedBy = &v
	}
}

// GetCreatedBy get created by
func (model *AuditedModel) GetCreatedBy() *string {
	return model.CreatedBy
}

// SetUpdatedBy set updated by
func (model *AuditedModel) SetUpdatedBy(updatedBy interface{}) {
	v := fmt.Sprintf("%v", updatedBy)
	model.UpdatedBy = &v
}

// GetUpdatedBy get updated by
func (model *AuditedModel) GetUpdatedBy() *string {
	return model.UpdatedBy
}

type auditableInterface interface {
	SetCreatedBy(createdBy interface{})
	GetCreatedBy() *string
	SetUpdatedBy(updatedBy interface{})
	GetUpdatedBy() *string
}
