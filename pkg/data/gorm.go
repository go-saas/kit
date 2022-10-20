package data

import (
	"gorm.io/gorm"
	"reflect"
)

func isModel[T any](db *gorm.DB) (t T, is bool) {
	if db.Statement.Schema.ModelType == nil {
		return
	}
	if db.Statement.Model != nil {
		t, is = db.Statement.Model.(T)
		if is {
			return
		}
	}
	_, is = reflect.New(db.Statement.Schema.ModelType).Interface().(T)
	return
}
