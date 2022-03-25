package gorm

import (
	sgorm "github.com/goxiaoy/go-saas/gorm"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/gorm"
)

func NewDbOpener() (sgorm.DbOpener, func()) {
	return sgorm.NewDbOpener(func(db *gorm.DB) *gorm.DB {
		RegisterCallbacks(db)
		if err := db.Use(otelgorm.NewPlugin()); err != nil {
			panic(err)
		}
		return db
	})
}
