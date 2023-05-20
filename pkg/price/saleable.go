package price

import (
	"context"
	"time"
)

// Saleable embed struct for saleable item
type Saleable struct {
	IsSaleable   bool
	SaleableFrom *time.Time
	SaleableTo   *time.Time

	Price Info `gorm:"embedded;embeddedPrefix:price_"`
}

// Info embed struct for holding price info
type Info struct {
	Default      Price `gorm:"embedded;embeddedPrefix:default_"`
	Discounted   Price `gorm:"embedded;embeddedPrefix:discounted_"`
	DiscountText string

	DenyMoreDiscounts bool
}

func (i Info) ToInfoPb(ctx context.Context) *InfoPb {
	return &InfoPb{
		Default:           i.Default.ToPricePb(ctx),
		Discounted:        i.Discounted.ToPricePb(ctx),
		DiscountText:      i.DiscountText,
		DenyMoreDiscounts: i.DenyMoreDiscounts,
	}
}

func NewInfoFromPb(i *InfoPb) (Info, error) {
	ret := Info{
		DiscountText:      i.DiscountText,
		DenyMoreDiscounts: i.DenyMoreDiscounts,
	}
	if i.Default != nil && len(i.Default.CurrencyCode) > 0 {
		p, err := NewPriceFromPb(i.Default)
		if err != nil {
			return ret, err
		}
		ret.Default = p
	}
	if i.Discounted != nil && len(i.Discounted.CurrencyCode) > 0 {
		p, err := NewPriceFromPb(i.Discounted)
		if err != nil {
			return ret, err
		}
		ret.Discounted = p
	}
	return ret, nil
}
