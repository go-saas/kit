package data

import (
	"github.com/go-saas/kit/pkg/authn"
	"gorm.io/gorm"
	"reflect"
)

func getCurrentUser(db *gorm.DB) (string, bool) {
	if u, ok := authn.FromUserContext(db.Statement.Context); ok {
		return u.GetId(), true
	}
	return "", false
}

func assignCreatedBy(db *gorm.DB) {
	if _, ok := isModel[auditableInterface](db); ok {
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
	if _, ok := isModel[auditableInterface](db); ok {
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
