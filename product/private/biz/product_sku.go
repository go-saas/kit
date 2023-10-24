package biz

import (
	kitgorm "github.com/go-saas/kit/pkg/gorm"
	"time"
)

// ProductSku sku
type ProductSku struct {
	kitgorm.UIDBase
	kitgorm.AuditedModel

	ProductId string
	Product   Product

	Title string

	MainPic   ProductMedia `gorm:"foreignKey:MainPicID"`
	MainPicID string
	Medias    []ProductMedia `gorm:"polymorphic:Owner;polymorphicValue:product_sku"`

	Prices []Price `gorm:"polymorphic:Owner;polymorphicValue:product_sku"`

	Stock []Stock `gorm:"polymorphic:Owner;polymorphicValue:product_sku"`

	Keyword []KeyWord `gorm:"polymorphic:Owner;polymorphicValue:product_sku;comment:商品关键字"`

	IsSaleable   bool
	SaleableFrom *time.Time
	SaleableTo   *time.Time
	Barcode      string `gorm:"comment:商品条码"`
}
