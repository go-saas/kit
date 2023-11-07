package biz

import (
	"github.com/go-saas/kit/pkg/data"
	"github.com/go-saas/kit/pkg/sortable"
)

type ProductMedia struct {
	ID string `gorm:"primaryKey;size:128"`

	OwnerID string
	// OwnerType product/product_sku
	OwnerType string

	Type      string
	MimeType  string
	Usage     string
	Name      string
	Reference string
	sortable.Embed
}

func NewProductMedia() *ProductMedia {
	return &ProductMedia{}
}

type ProductMediaRepo interface {
	data.Repo[ProductMedia, string, interface{}]
}
