package biz

import kitgorm "github.com/go-saas/kit/pkg/gorm"

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
