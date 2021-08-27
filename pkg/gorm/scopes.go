package gorm

import "gorm.io/gorm"

// WhereUserId append 'user_id' field filter
func WhereUserId(id interface{}) func(db *gorm.DB) *gorm.DB  {
	return func(db *gorm.DB) *gorm.DB{
		return db.Where("user_id = ?", id)
	}
}

// WhereIf append where if conditional return true
func WhereIf(conditional func() bool,query interface{}, args ...interface{} )func(db *gorm.DB) *gorm.DB   {
	return func(db *gorm.DB) *gorm.DB {
		v:= conditional()
		if v{
			return db.Where(query,args...)
		}
		return db
	}
}

// OrIf append or if conditional return true
func OrIf (conditional func() bool,query interface{}, args ...interface{} )func(db *gorm.DB) *gorm.DB   {
	return func(db *gorm.DB) *gorm.DB {
		v:= conditional()
		if v{
			return db.Or(query,args...)
		}
		return db
	}
}