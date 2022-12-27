package gorm

import (
	"fmt"
	"github.com/go-saas/kit/pkg/authn"
	"gorm.io/gorm"
	"reflect"
	"time"
)

// AuditedModel make Model Auditable, embed `audited.AuditedModel` into your model as anonymous field to make the model auditable
//
//	type User struct {
//	  gorm.AuditedModel
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

func getCurrentUser(db *gorm.DB) (string, bool) {
	if u, ok := authn.FromUserContext(db.Statement.Context); ok {
		return u.GetId(), true
	}
	return "", false
}

func assignCreatedBy(db *gorm.DB) {
	if _, ok := IsModel[auditableInterface](db); ok {
		if user, ok := getCurrentUser(db); ok {
			f := db.Statement.Schema.FieldsByName["CreatedBy"]
			switch db.Statement.ReflectValue.Kind() {
			case reflect.Slice, reflect.Array:
				for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
					f.Set(db.Statement.Context, db.Statement.ReflectValue.Index(i), &user)
				}
			case reflect.Struct:
				f.Set(db.Statement.Context, db.Statement.ReflectValue, &user)
			}
		}
	}
}

func assignUpdatedBy(db *gorm.DB) {
	if _, ok := IsModel[auditableInterface](db); ok {
		if user, ok := getCurrentUser(db); ok {
			f := db.Statement.Schema.FieldsByName["UpdatedBy"]
			switch db.Statement.ReflectValue.Kind() {
			case reflect.Slice, reflect.Array:
				for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
					f.Set(db.Statement.Context, db.Statement.ReflectValue.Index(i), &user)
				}
			case reflect.Struct:
				f.Set(db.Statement.Context, db.Statement.ReflectValue, &user)
			}
		}
	}
}

// RegisterAuditCallbacks register callback into GORM DB
func RegisterAuditCallbacks(db *gorm.DB) {
	callback := db.Callback()
	if callback.Create().Get("audited:assign_created_by") == nil {
		callback.Create().Before("gorm:before_create").Register("audited:assign_created_by", assignCreatedBy)
	}
	if callback.Update().Get("audited:assign_updated_by") == nil {
		callback.Update().Before("gorm:before_update").Register("audited:assign_updated_by", assignUpdatedBy)
	}
}
