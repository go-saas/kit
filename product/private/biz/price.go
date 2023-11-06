package biz

import (
	"github.com/go-saas/kit/pkg/data"
	kitgorm "github.com/go-saas/kit/pkg/gorm"
	v1 "github.com/go-saas/kit/product/api/price/v1"
	"github.com/go-saas/saas/gorm"
)

// Describes how to compute the price per period. Either `per_unit` or `tiered`. `per_unit` indicates that the fixed amount (specified in `unit_amount` or `unit_amount_decimal`) will be charged per unit in `quantity` (for prices with `usage_type=licensed`), or per unit of total usage (for prices with `usage_type=metered`). `tiered` indicates that the unit pricing will be computed using a tiering strategy as defined using the `tiers` and `tiers_mode` attributes.
type PriceBillingScheme string

// List of values that PriceBillingScheme can take
const (
	PriceBillingSchemePerUnit PriceBillingScheme = "per_unit"
	PriceBillingSchemeTiered  PriceBillingScheme = "tiered"
)

// Specifies a usage aggregation strategy for prices of `usage_type=metered`. Allowed values are `sum` for summing up all usage during a period, `last_during_period` for using the last usage record reported within a period, `last_ever` for using the last usage record ever (across period bounds) or `max` which uses the usage record with the maximum reported usage during a period. Defaults to `sum`.
type PriceRecurringAggregateUsage string

// List of values that PriceRecurringAggregateUsage can take
const (
	PriceRecurringAggregateUsageLastDuringPeriod PriceRecurringAggregateUsage = "last_during_period"
	PriceRecurringAggregateUsageLastEver         PriceRecurringAggregateUsage = "last_ever"
	PriceRecurringAggregateUsageMax              PriceRecurringAggregateUsage = "max"
	PriceRecurringAggregateUsageSum              PriceRecurringAggregateUsage = "sum"
)

// PriceRecurringInterval is the list of allowed values for a price's recurring interval.
type PriceRecurringInterval string

// List of values that PriceRecurringInterval can take.
const (
	PriceRecurringIntervalDay   PriceRecurringInterval = "day"
	PriceRecurringIntervalWeek  PriceRecurringInterval = "week"
	PriceRecurringIntervalMonth PriceRecurringInterval = "month"
	PriceRecurringIntervalYear  PriceRecurringInterval = "year"
)

// Configures how the quantity per period should be determined. Can be either `metered` or `licensed`. `licensed` automatically bills the `quantity` set when adding it to a subscription. `metered` aggregates the total usage based on usage records. Defaults to `licensed`.
type PriceRecurringUsageType string

// List of values that PriceRecurringUsageType can take
const (
	PriceRecurringUsageTypeLicensed PriceRecurringUsageType = "licensed"
	PriceRecurringUsageTypeMetered  PriceRecurringUsageType = "metered"
)

// Defines if the tiering price should be `graduated` or `volume` based. In `volume`-based tiering, the maximum quantity within a period determines the per unit price. In `graduated` tiering, pricing can change as the quantity grows.
type PriceTiersMode string

// List of values that PriceTiersMode can take
const (
	PriceTiersModeGraduated PriceTiersMode = "graduated"
	PriceTiersModeVolume    PriceTiersMode = "volume"
)

// After division, either round the result `up` or `down`.
type PriceTransformQuantityRound string

// List of values that PriceTransformQuantityRound can take
const (
	PriceTransformQuantityRoundDown PriceTransformQuantityRound = "down"
	PriceTransformQuantityRoundUp   PriceTransformQuantityRound = "up"
)

// PriceType is the list of allowed values for a price's type.
type PriceType string

// List of values that PriceType can take.
const (
	PriceTypeOneTime   PriceType = "one_time"
	PriceTypeRecurring PriceType = "recurring"
)

type Price struct {
	kitgorm.UIDBase
	kitgorm.AuditedModel
	gorm.MultiTenancy

	OwnerID string
	// OwnerType product/product_sku
	OwnerType string

	DefaultAmount     int64
	DiscountedAmount  *int64
	CurrencyCode      string
	DiscountText      string
	DenyMoreDiscounts bool

	BillingScheme   PriceBillingScheme
	CurrencyOptions []PriceCurrencyOption `gorm:"foreignKey:PriceId"`
	Recurring       *PriceRecurring       `gorm:"foreignKey:PriceId"`

	Tiers     []PriceTier
	TiersMode PriceTiersMode `json:"tiers_mode"`

	TransformQuantity PriceTransformQuantity `gorm:"embedded"`

	ProductID string

	Type PriceType

	StripePriceId *string
}

func NewPrice(productID string) *Price {
	return &Price{ProductID: productID}
}

func (p *Price) GetNeedPayAmount() int64 {
	if p.DiscountedAmount != nil {
		return *p.DiscountedAmount
	}
	return p.DefaultAmount
}

type PriceCurrencyOption struct {
	kitgorm.UIDBase
	PriceId string

	DefaultAmount     int64
	DiscountedAmount  *int64
	CurrencyCode      string
	DiscountText      string
	DenyMoreDiscounts bool

	Tiers []PriceCurrencyOptionTier `gorm:"foreignKey:PriceCurrencyOptionId"`
}

func (p *PriceCurrencyOption) GetNeedPayAmount() int64 {
	if p.DiscountedAmount != nil {
		return *p.DiscountedAmount
	}
	return p.DefaultAmount
}

type PriceRecurring struct {
	kitgorm.UIDBase
	PriceId string

	AggregateUsage  PriceRecurringAggregateUsage
	Interval        PriceRecurringInterval
	IntervalCount   int64
	TrialPeriodDays int64
	UsageType       PriceRecurringUsageType
}

type PriceTier struct {
	kitgorm.UIDBase
	PriceId string

	// Price for the entire tier.
	FlatAmount int64 `json:"flat_amount"`
	// Per unit price for units relevant to the tier.
	UnitAmount int64 `json:"unit_amount"`
	// Up to and including to this quantity will be contained in the tier.
	UpTo int64 `json:"up_to"`
}

// Apply a transformation to the reported usage or set quantity before computing the amount billed. Cannot be combined with `tiers`.
type PriceTransformQuantity struct {
	// Divide usage by this number.
	DivideBy int64 `json:"divide_by"`
	// After division, either round the result `up` or `down`.
	Round PriceTransformQuantityRound `json:"round"`
}

// Each element represents a pricing tier. This parameter requires `billing_scheme` to be set to `tiered`. See also the documentation for `billing_scheme`.
type PriceCurrencyOptionTier struct {
	kitgorm.UIDBase
	PriceCurrencyOptionId string

	// Price for the entire tier.
	FlatAmount int64 `json:"flat_amount"`
	// Per unit price for units relevant to the tier.
	UnitAmount int64 `json:"unit_amount"`
	// Up to and including to this quantity will be contained in the tier.
	UpTo int64 `json:"up_to"`
}

type PriceRepo interface {
	data.Repo[Price, string, *v1.ListPriceRequest]
}
