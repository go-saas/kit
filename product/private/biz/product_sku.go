package biz

import (
	kitgorm "github.com/go-saas/kit/pkg/gorm"
	"github.com/go-saas/kit/pkg/price"
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

	Price price.Info `gorm:"embedded;embeddedPrefix:price_"`

	Stock []Stock `gorm:"polymorphic:Owner;polymorphicValue:product_sku"`

	Barcode string `gorm:"comment:商品条码"`
}
