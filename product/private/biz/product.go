package biz

import (
	"context"
	"github.com/go-saas/kit/pkg/data"
	kitgorm "github.com/go-saas/kit/pkg/gorm"
	"github.com/go-saas/kit/pkg/sortable"
	v1 "github.com/go-saas/kit/product/api/product/v1"
	sgorm "github.com/go-saas/saas/gorm"
	"gorm.io/gorm"

	concurrency "github.com/goxiaoy/gorm-concurrency/v2"
	"time"
)

// Product SPU
type Product struct {
	kitgorm.UIDBase
	kitgorm.AuditedModel
	concurrency.HasVersion
	sgorm.MultiTenancy
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Title     string       `gorm:"comment:商品标题"`
	ShortDesc string       `gorm:"comment:商品简述"`
	Desc      string       `gorm:"comment:商品描述"`
	Content   data.JSONMap `gorm:"comment:描述页面内容"`

	MainPic *ProductMedia  `gorm:"polymorphic:Owner;polymorphicValue:product"`
	Medias  []ProductMedia `gorm:"polymorphic:Owner;polymorphicValue:product"`

	Badges []Badge `gorm:"foreignKey:ProductId;comment:商品徽章"`

	VisibleFrom *time.Time
	VisibleTo   *time.Time

	IsNew bool

	Categories []ProductCategory `gorm:"many2many:product_2_category"`

	MainCategoryKey *string
	MainCategory    *ProductCategory `gorm:"foreignKey:MainCategoryKey;comment:商品主要分类"`

	Keywords []Keyword `gorm:"polymorphic:Owner;polymorphicValue:product;comment:商品关键字"`

	Model string `gorm:"comment:商品型号"`

	BrandId *string
	Brand   *Brand

	IsGiveaway bool `gorm:"comment:是否赠品"`

	Attributes []ProductAttribute

	MultiSku bool `gorm:"comment:是否多SKU产品,只能创建时修改"`

	// CampaignRules
	CampaignRules []CampaignRule `gorm:"polymorphic:Owner;polymorphicValue:product"`

	NeedShipping bool `gorm:"comment:是否需要邮寄"`

	Stocks []Stock `gorm:"polymorphic:Owner;polymorphicValue:product"`

	Sku []ProductSku `gorm:"foreignKey:ProductId"`

	Prices []Price `gorm:"polymorphic:Owner;polymorphicValue:product"`

	//Product Basic
	SaleableFrom *time.Time
	SaleableTo   *time.Time

	Barcode string `gorm:"comment:商品条码"`

	ManageInfo ProductManageInfo `gorm:"embedded"`
	SyncLinks  []ProductSyncLink `gorm:"foreignKey:ProductId"`

	Active bool

	Internal bool
}

type ProductManageInfo struct {
	Managed      bool
	ManagedBy    string
	LastSyncTime *time.Time
}

type Badge struct {
	kitgorm.UIDBase
	ProductId string

	Code  string
	Label string
}

type Keyword struct {
	kitgorm.UIDBase
	OwnerID string
	// OwnerType product/product_sku
	OwnerType string

	Text  string
	Refer string
}

type ProductRepo interface {
	data.Repo[Product, string, *v1.ListProductRequest]
	GetSyncLinks(ctx context.Context, product *Product) ([]ProductSyncLink, error)
	UpdateSyncLink(ctx context.Context, product *Product, syncLink *ProductSyncLink) error
}

type CampaignRule struct {
	kitgorm.UIDBase
	OwnerID string
	// OwnerType product or product_sku
	OwnerType string
	Rule      string
	Extra     data.JSONMap
}

// PriceContext defines the scope in which the price was calculated
type PriceContext struct {
}

// ProductAttribute TODO how to add custom attribute
type ProductAttribute struct {
	kitgorm.UIDBase
	ProductId string

	Title string
	sortable.Embed
}

// Stock holds data with product availability info
type Stock struct {
	kitgorm.UIDBase

	OwnerID string
	// OwnerType product/product_sku
	OwnerType    string
	InStock      bool
	Level        string
	Amount       int
	DeliveryCode string
}

// Stock Level values
const (
	StockLevelOutOfStock = "out"
	StockLevelInStock    = "in"
	StockLevelLowStock   = "low"
)
