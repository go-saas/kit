package biz

import (
	v1 "github.com/goxiaoy/go-saas-kit/payment/api/order/v1"
	"github.com/goxiaoy/go-saas-kit/pkg/data"
	"github.com/goxiaoy/go-saas-kit/pkg/gorm"
	sgorm "github.com/goxiaoy/go-saas/gorm"
)

type Order struct {
	gorm.UIDBase
	gorm.AuditedModel
	sgorm.MultiTenancy
	Name string
}

type OrderRepo interface {
	data.Repo[Order, string, v1.ListOrderRequest]
}
