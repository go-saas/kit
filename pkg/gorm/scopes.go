package gorm

import "gorm.io/gorm"

func WhereUserId(db *gorm.DB, id interface{}) *gorm.DB {
	return db.Where("user_id = ?", id)
}
