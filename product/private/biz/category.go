package biz

import (
	"context"
	"github.com/go-saas/kit/pkg/data"
	kitgorm "github.com/go-saas/kit/pkg/gorm"
	v1 "github.com/go-saas/kit/product/api/category/v1"
)

// ProductCategory represents some Teaser infos for ProductCategory
type ProductCategory struct {
	// Key the identifier of the ProductCategory
	Key string `gorm:"primaryKey;size:128"`

	kitgorm.AuditedModel
	// The Path (root to leaf) for this ProductCategory - separated by "/"
	Path string
	// Name is the speaking name of the category
	Name     string
	ParentID *string

	// Parent is an optional link to parent teaser
	Parent *ProductCategory `gorm:"foreignKey:ParentID;references:key"`
}

type ProductCategoryRepo interface {
	data.Repo[ProductCategory, string, *v1.ListProductCategoryRequest]
	FindAllChildren(ctx context.Context, entity *ProductCategory) ([]*ProductCategory, error)
	FindByKeys(ctx context.Context, cKeys []string) ([]ProductCategory, error)
}
