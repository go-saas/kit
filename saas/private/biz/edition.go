package biz

import (
	"github.com/go-saas/kit/pkg/data"
	"github.com/go-saas/kit/pkg/gorm"
)

type Edition struct {
	gorm.UIDBase
	Name        string
	DisplayName string
	Features    []EditionFeature `gorm:"foreignKey:EditionId"`
}

type EditionFeature struct {
	gorm.UIDBase
	EditionId string
	Key       string     `gorm:"column:key;primary_key;size:100;"`
	Value     data.Value `gorm:"embedded"`
}
