package biz

import (
	"context"
	"github.com/go-saas/kit/pkg/data"
	"github.com/go-saas/kit/pkg/gorm"
	"github.com/go-saas/kit/pkg/sortable"
	v1 "github.com/go-saas/kit/saas/api/plan/v1"
)

type Plan struct {
	gorm.AuditedModel
	sortable.Embed
	Key         string `gorm:"primaryKey;size:128"`
	DisplayName string
	Active      bool

	//ProductId linked to product service
	ProductId string `gorm:"size:200;index"`

	Features []PlanFeature `gorm:"foreignKey:PlanId"`
}

type PlanFeature struct {
	gorm.UIDBase
	PlanId string
	Key    string     `gorm:"column:key;primary_key;size:100;"`
	Value  data.Value `gorm:"embedded"`
}

func NewPlan(key, displayName, productId string, sort int) *Plan {
	res := &Plan{
		Key:         key,
		DisplayName: displayName,
		ProductId:   productId,
		Active:      false,
	}
	res.Sort = sort
	return res
}

type PlanRepo interface {
	data.Repo[Plan, string, *v1.ListPlanRequest]
	FindByProductId(ctx context.Context, productID string) (*Plan, error)
}
