package gorm

import "gorm.io/gorm"

// WhereUserId append 'user_id' field filter
func WhereUserId(id string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", id)
	}
}
