package biz

import (
	v1 "github.com/go-saas/kit/payment/api/subscription/v1"
	"github.com/go-saas/kit/pkg/data"
	kitgorm "github.com/go-saas/kit/pkg/gorm"
	sgorm "github.com/go-saas/saas/gorm"
	"time"
)

type SubscriptionStatus string

const (
	SubscriptionStatusActive            SubscriptionStatus = "active"
	SubscriptionStatusCanceled          SubscriptionStatus = "canceled"
	SubscriptionStatusIncomplete        SubscriptionStatus = "incomplete"
	SubscriptionStatusIncompleteExpired SubscriptionStatus = "incomplete_expired"
	SubscriptionStatusPastDue           SubscriptionStatus = "past_due"
	SubscriptionStatusPaused            SubscriptionStatus = "paused"
	SubscriptionStatusTrialing          SubscriptionStatus = "trialing"
	SubscriptionStatusUnpaid            SubscriptionStatus = "unpaid"
)

type Subscription struct {
	kitgorm.UIDBase
	kitgorm.AuditedModel
	sgorm.MultiTenancy

	UserId string `gorm:"index;size:128"`

	CancelAtPeriodEnd bool
	CurrencyCode      string

	CurrentPeriodStart *time.Time
	CurrentPeriodEnd   *time.Time

	Status SubscriptionStatus `gorm:"index;size:128"`

	Provider    string `gorm:"index;size:128"`
	ProviderKey string `gorm:"index;size:128"`

	Items []SubscriptionItem `gorm:"foreignKey:SubscriptionID"`
}

type SubscriptionItem struct {
	kitgorm.UIDBase
	kitgorm.AuditedModel
	SubscriptionID string
	//PriceID linked with product.Prices
	PriceID        string
	ProductOrSkuID string
	ProductType    string
}

type SubscriptionListPrams interface {
	GetSearch() string
	GetFilter() *v1.SubscriptionFilter
}

type SubscriptionRepo interface {
	data.Repo[Subscription, string, SubscriptionListPrams]
}
