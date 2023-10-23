package biz

import (
	"github.com/go-saas/kit/pkg/data"
	"github.com/go-saas/kit/pkg/gorm"
	v1 "github.com/go-saas/kit/saas/api/plan/v1"
)

type Plan struct {
	gorm.AuditedModel
	Key         string `gorm:"primaryKey;size:128"`
	DisplayName string
	Active      bool
	Features    []PlanFeature `gorm:"foreignKey:PlanId"`
}

type PlanFeature struct {
	gorm.UIDBase
	PlanId string
	Key    string     `gorm:"column:key;primary_key;size:100;"`
	Value  data.Value `gorm:"embedded"`
}

type PlanRepo interface {
	data.Repo[Plan, string, *v1.ListPlanRequest]
}
