package biz

import (
	"context"
	"github.com/cockroachdb/apd/v3"
	v1 "github.com/go-saas/kit/order/api/order/v1"
	"github.com/go-saas/kit/pkg/data"
	kitgorm "github.com/go-saas/kit/pkg/gorm"
	"github.com/go-saas/kit/pkg/price"
	"github.com/go-saas/lbs"
	gorm2 "github.com/go-saas/saas/gorm"
	concurrency "github.com/goxiaoy/gorm-concurrency/v2"
	"github.com/lithammer/shortuuid/v3"
	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
	"time"
)

const (
	OrderStatusUnpaid    string = "UNPAID"
	OrderStatusPaid      string = "PAID"
	OrderStatusRefunding string = "REFUNDING"
	OrderStatusRefunded  string = "REFUNDED"
	OrderStatusExpired   string = "EXPIRED"
)

const (
	OrderFlowTypePay           string = "PAID"
	OrderFlowTypeRequestPay    string = "REQUEST_PAY"
	OrderFlowTypeRequestRefund string = "REQUEST_REFUND"
	OrderFlowTypeRefund        string = "REFUND"
)

type Order struct {
	ID string `gorm:"type:char(36)" json:"id"`
	kitgorm.AuditedModel
	concurrency.HasVersion
	gorm2.MultiTenancy
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Status string

	CurrencyCode string

	TotalPriceAmount        int64
	TotalPriceTaxAmount     int64
	TotalPriceInclTaxAmount int64

	DiscountAmount int64

	OriginalPriceAmount int64 `gorm:"comment:TotalPrice+Discount"`

	PaidPriceAmount int64

	PaidTime *time.Time

	PayBefore *time.Time

	PayProvider string
	PayMethod   string

	ShippingAddr lbs.AddressEntity `gorm:"embedded;embeddedPrefix:shipping_addr_"`
	BillingAddr  lbs.AddressEntity `gorm:"embedded;embeddedPrefix:billing_addr_"`

	CustomerID string `gorm:"size:200;index:,;comment:一般等于用户ID"`

	Extra data.JSONMap

	Items []OrderItem `gorm:"foreignKey:OrderID;references:ID"`

	PaymentProviders []OrderPaymentProvider `gorm:"foreignKey:OrderID;references:ID"`
}

type OrderItem struct {
	ID string `gorm:"type:char(36)" json:"id"`

	Qty int64

	CurrencyCode string

	PriceAmount        int64
	PriceTaxAmount     int64
	PriceInclTaxAmount int64

	OriginalPriceAmount int64

	RowTotalAmount         int64
	RowTotalTaxAmount      int64
	RowTotalInclTaxAmount  int64 `gorm:"comment:SinglePriceInclTax*Qty,RowTotal+RowTax"`
	OriginalRowTotalAmount int64

	RowDiscountAmount int64

	Product OrderProduct `gorm:"embedded;"`

	OrderID    string
	IsGiveaway bool `gorm:"comment:是否赠品"`

	BizPayload data.JSONMap
}

type OrderPaymentProvider struct {
	kitgorm.UIDBase
	gorm2.MultiTenancy
	OrderID     string
	Provider    string
	ProviderKey string
}

func NewOrder(currencyCode string, taxRate apd.Decimal, items []OrderItem) (*Order, error) {
	i := &Order{Items: items, Status: OrderStatusUnpaid, CurrencyCode: currencyCode}

	var totalPrice int64
	var originalPrice int64
	var err error
	for i, item := range i.Items {
		if i == 0 {
			totalPrice = item.RowTotalAmount
			rowOriginalPrice := item.RowTotalAmount * item.Qty

			originalPrice = rowOriginalPrice
		} else {
			totalPrice = totalPrice + item.RowTotalAmount
			rowOriginalPrice := item.OriginalPriceAmount

			originalPrice = originalPrice + rowOriginalPrice
		}
	}
	i.TotalPriceAmount = totalPrice

	i.TotalPriceTaxAmount, i.TotalPriceInclTaxAmount, err = calWithTax(i.TotalPriceAmount, i.CurrencyCode, taxRate)
	if err != nil {
		return nil, err
	}
	i.OriginalPriceAmount = originalPrice

	i.DiscountAmount = originalPrice - totalPrice
	return i, nil
}

func (i *OrderItem) BeforeCreate(tx *gorm.DB) error {
	if len(i.ID) == 0 {
		i.ID = shortuuid.New()
	}
	return nil
}

func NewOrderItemFromPriceAndOriginalPrice(
	currencyCode string,
	product OrderProduct,
	qty int64,
	taxRate apd.Decimal,
	priceAmount int64,
	originalPriceAmount int64,
	isGiveaway bool,
	bizPayload data.JSONMap,
) (*OrderItem, error) {
	i := &OrderItem{
		CurrencyCode:        currencyCode,
		Qty:                 qty,
		PriceAmount:         priceAmount,
		OriginalPriceAmount: originalPriceAmount,
		Product:             product,
		IsGiveaway:          isGiveaway,
		BizPayload:          bizPayload,
	}

	singleDiscount := i.OriginalPriceAmount - i.PriceAmount

	i.RowTotalAmount = i.PriceAmount * i.Qty

	i.RowDiscountAmount = singleDiscount * i.Qty
	i.OriginalRowTotalAmount = i.OriginalPriceAmount * i.Qty

	return i, i.Cal(taxRate)
}

func (i *OrderItem) Cal(taxRate apd.Decimal) (err error) {

	i.PriceTaxAmount, i.PriceInclTaxAmount, err = calWithTax(i.PriceAmount, i.CurrencyCode, taxRate)

	if err != nil {
		return err
	}
	i.RowTotalTaxAmount, i.RowTotalInclTaxAmount, err = calWithTax(i.RowTotalAmount, i.CurrencyCode, taxRate)

	if err != nil {
		return err
	}
	return
}

func calWithTax(p int64, currencyCode string, taxRate apd.Decimal) (taxAmount int64, priceInclTaxAmount int64, err error) {
	pp := price.MustNewFromInt64(p, currencyCode).ToCurrency()

	taxx, err := pp.Mul(taxRate.String())
	if err != nil {
		return
	}
	taxAmount, err = taxx.Int64()
	if err != nil {
		return
	}
	priceInclTaxx, err := pp.Add(taxx)
	if err != nil {
		return
	}
	priceInclTaxAmount, err = priceInclTaxx.Int64()
	return
}

type OrderProduct struct {
	ProductName     string
	ProductMainPic  string
	ProductID       *string `gorm:"type:char(36);index:,"`
	ProductVersion  string  `gorm:"type:char(36);index:,"`
	ProductType     string  `gorm:"size:128;index:,"`
	ProductSkuID    *string `gorm:"type:char(36);index:,"`
	ProductSkuTitle string
	PriceID         *string `gorm:"type:char(36);index:,"`
}

type OrderRepo interface {
	data.Repo[Order, string, *v1.ListOrderRequest]
	FindByPaymentProvider(ctx context.Context, provider, providerKey string) (*Order, error)
	UpsertPaymentProvider(ctx context.Context, order *Order, provider *OrderPaymentProvider) error
	ListPaymentProviders(ctx context.Context, order *Order) ([]OrderPaymentProvider, error)
}

func (u *Order) BeforeCreate(tx *gorm.DB) error {
	if len(u.ID) == 0 {
		u.ID = ksuid.New().String()
	}
	return nil
}

func (u *Order) ChangeToPaid(payProvider string, paymethod string, paidPriceAmount int64, paidTime *time.Time) {
	u.PayProvider = payProvider
	u.PayMethod = paymethod
	u.PaidPriceAmount = paidPriceAmount
	u.PaidTime = paidTime
	u.Status = OrderStatusPaid
}

func (u *Order) RequestFund(payway string, refundPriceAmount int64, data map[string]interface{}) {
	u.Status = OrderStatusRefunding
	//TODO
}

func (u *Order) ChangeToRefunded(payway string, refundedPriceAmount int64, data map[string]interface{}) {
	u.Status = OrderStatusRefunded
	//TODO
}
