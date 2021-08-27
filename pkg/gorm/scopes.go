package gorm

import "gorm.io/gorm"

func WhereUserId(id interface{}) func(db *gorm.DB) *gorm.DB  {
	return func(db *gorm.DB) *gorm.DB{
		return db.Where("user_id = ?", id)
	}
}