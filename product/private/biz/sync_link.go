package biz

import "time"

type ProductSyncLink struct {
	ProductId    string `gorm:"primary_key"`
	ProviderName string `gorm:"primary_key"`
	ProviderId   string `gorm:"index"`
	LastSyncTime *time.Time
}
