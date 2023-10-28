package price

import (
	"context"
	"github.com/bojanz/currency"
	"github.com/go-saas/kit/pkg/localize"
)

// Price database friendly price struct
type Price struct {
	//Amount minor units (e.g. cents). for some database(like sqlite), decimal is stored as float, it may cause problems when calculating
	Amount       int64 `json:"amount"`
	CurrencyCode string
}

// NewPriceFromCurrency convert currency.Amount into database friendly Price
func NewPriceFromCurrency(a currency.Amount) (Price, error) {
	v, err := a.Int64()
	if err != nil {
		return Price{}, err
	}
	return Price{Amount: v, CurrencyCode: a.CurrencyCode()}, nil
}

func NewPrice(n, currencyCode string) (p Price, err error) {
	amount, err := currency.NewAmount(n, currencyCode)
	if err != nil {
		return
	}
	return NewPriceFromCurrency(amount)
}

func MustNew(n, currencyCode string) Price {
	p, err := NewPrice(n, currencyCode)
	if err != nil {
		panic(err)
	}
	return p
}

func NewPriceFromInt64(n int64, currencyCode string) (p Price, err error) {
	amount, err := currency.NewAmountFromInt64(n, currencyCode)
	if err != nil {
		return
	}
	return NewPriceFromCurrency(amount)
}

func MustNewFromInt64(n int64, currencyCode string) Price {
	p, err := NewPriceFromInt64(n, currencyCode)
	if err != nil {
		panic(err)
	}
	return p
}

func NewPriceFromPb(a *PricePb) (Price, error) {
	return NewPriceFromInt64(a.Amount, a.CurrencyCode)
}

func (p Price) ToCurrency() currency.Amount {
	v, err := currency.NewAmountFromInt64(p.Amount, p.CurrencyCode)
	if err != nil {
		panic(err)
	}
	return v
}

func (p Price) IsEmpty() bool {
	return p.Amount == 0 && p.CurrencyCode == ""
}

func (p Price) ToPricePb(ctx context.Context) *PricePb {
	if len(p.CurrencyCode) == 0 {
		return nil
	}
	a := p.ToCurrency()
	d, _ := currency.GetDigits(p.CurrencyCode)
	tags := localize.LanguageTags(ctx)
	var s string
	for _, tag := range tags {
		locale := currency.NewLocale(tag.String())
		if symbol, ok := currency.GetSymbol(p.CurrencyCode, locale); ok {
			s = symbol
			break
		}
	}
	if s == "" {
		s = p.CurrencyCode
	}
	return &PricePb{
		Amount:        p.Amount,
		AmountDecimal: a.Number(),
		CurrencyCode:  p.CurrencyCode,
		Text:          s + a.Number(),
		Digits:        int32(d),
	}
}
