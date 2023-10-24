package biz

import (
	"github.com/go-saas/kit/pkg/data"
	"github.com/go-saas/kit/pkg/gorm"
	"github.com/go-saas/kit/pkg/sortable"
	v1 "github.com/go-saas/kit/saas/api/plan/v1"
)

type Plan struct {
	gorm.AuditedModel
	*sortable.Embed
	Key         string `gorm:"primaryKey;size:128"`
	DisplayName string
	Active      bool

	//ProductId linked to product service
	ProductId string

	Features []PlanFeature `gorm:"foreignKey:PlanId"`
}

type PlanFeature struct {
	gorm.UIDBase
	PlanId string
	Key    string     `gorm:"column:key;primary_key;size:100;"`
	Value  data.Value `gorm:"embedded"`
}

func NewPlan(key, displayName string) *Plan {
	return &Plan{
		Embed:       &sortable.Embed{},
		Key:         key,
		DisplayName: displayName,
		Active:      false,
	}
}

type PlanRepo interface {
	data.Repo[Plan, string, *v1.ListPlanRequest]
}
