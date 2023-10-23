package biz

import (
	"fmt"
	"github.com/bojanz/currency"
	"github.com/cockroachdb/apd/v3"
	v1 "github.com/go-saas/kit/order/api/order/v1"
	"github.com/go-saas/kit/pkg/data"
	kitgorm "github.com/go-saas/kit/pkg/gorm"
	"github.com/go-saas/kit/pkg/price"
	"github.com/go-saas/lbs"
	concurrency "github.com/goxiaoy/gorm-concurrency"
	"github.com/lithammer/shortuuid/v3"
	"github.com/samber/lo"
	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type Order struct {
	ID string `gorm:"type:char(36)" json:"id"`
	kitgorm.AuditedModel
	concurrency.Version

	Status string

	TotalPrice        price.Price `gorm:"embedded;embeddedPrefix:total_price_"`
	TotalPriceInclTax price.Price `gorm:"embedded;embeddedPrefix:total_price_incl_tax_"`

	Discount price.Price `gorm:"embedded;embeddedPrefix:discount_"`

	OriginalPrice        price.Price `gorm:"embedded;embeddedPrefix:original_price_;comment:TotalPrice+Discount"`
	OriginalPriceInclTax price.Price `gorm:"embedded;embeddedPrefix:original_price_incl_tax_"`

	PaidPrice price.Price `gorm:"embedded;embeddedPrefix:paid_price_"`
	PaidTime  *time.Time

	PayBefore *time.Time

	PayWay    string
	PayMethod string

	FlowData []OrderFlowData `gorm:"foreignKey:OrderID;references:ID"`

	ShippingAddr lbs.AddressEntity `gorm:"embedded;embeddedPrefix:shipping_addr_"`
	BillingAddr  lbs.AddressEntity `gorm:"embedded;embeddedPrefix:billing_addr_"`

	CustomerID string `gorm:"type:char(36);index:,;comment:一般等于用户ID"`

	Extra data.JSONMap

	Items []OrderItem `gorm:"foreignKey:OrderID;references:ID"`
}

type OrderFlowData struct {
	kitgorm.UIDBase
	OrderID     string
	PayWay      string
	FlowType    string
	Price       price.Price `gorm:"embedded"`
	InitialTime time.Time
	Data        data.JSONMap
}

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

var (
	ErrOrderItemsRequired = fmt.Errorf("order items required")
)

func NewOrder(taxRate apd.Decimal, items []OrderItem) (*Order, error) {
	if len(items) == 0 {
		return nil, ErrOrderItemsRequired
	}
	i := &Order{Items: items, Status: OrderStatusUnpaid}

	var totalPrice currency.Amount
	var originalPrice currency.Amount
	var err error
	for i, item := range i.Items {
		if i == 0 {
			totalPrice = item.RowTotal.ToCurrency()
			rowOriginalPrice, err := item.OriginalPrice.ToCurrency().Mul(strconv.FormatInt(int64(item.Qty), 10))
			if err != nil {
				return nil, err
			}
			originalPrice = rowOriginalPrice
		} else {
			totalPrice, err = totalPrice.Add(item.RowTotal.ToCurrency())
			if err != nil {
				return nil, err
			}
			rowOriginalPrice, err := item.OriginalPrice.ToCurrency().Mul(strconv.FormatInt(int64(item.Qty), 10))
			if err != nil {
				return nil, err
			}
			originalPrice, err = originalPrice.Add(rowOriginalPrice)
			if err != nil {
				return nil, err
			}
		}
	}
	i.TotalPrice, err = price.NewPriceFromCurrency(totalPrice)
	if err != nil {
		return nil, err
	}
	_, i.TotalPriceInclTax, err = calWithTax(i.TotalPrice, taxRate)
	if err != nil {
		return nil, err
	}
	i.OriginalPrice, err = price.NewPriceFromCurrency(originalPrice)
	if err != nil {
		return nil, err
	}
	_, i.OriginalPriceInclTax, err = calWithTax(i.OriginalPrice, taxRate)
	if err != nil {
		return nil, err
	}
	discount, err := originalPrice.Sub(totalPrice)
	if err != nil {
		return nil, err
	}
	i.Discount, err = price.NewPriceFromCurrency(discount)
	if err != nil {
		return nil, err
	}
	return i, nil
}

type OrderItem struct {
	ID string `gorm:"type:char(36)" json:"id"`

	Qty int32

	Price        price.Price `gorm:"embedded;embeddedPrefix:price_"`
	Tax          price.Price `gorm:"embedded;embeddedPrefix:tax_"`
	PriceInclTax price.Price `gorm:"embedded;embeddedPrefix:price_incl_tax_;comment:SinglePrice+SingleTax"`

	RowTotal        price.Price `gorm:"embedded;embeddedPrefix:row_total_"`
	RowTotalTax     price.Price `gorm:"embedded;embeddedPrefix:row_total_tax"`
	RowTotalInclTax price.Price `gorm:"embedded;embeddedPrefix:row_total_incl_tax_;comment:SinglePriceInclTax*Qty,RowTotal+RowTax"`

	OriginalPrice        price.Price `gorm:"embedded;embeddedPrefix:original_price_"`
	OriginalPriceTax     price.Price `gorm:"embedded;embeddedPrefix:original_price_tax_"`
	OriginalPriceInclTax price.Price `gorm:"embedded;embeddedPrefix:original_price_incl_tax_"`

	RowDiscount price.Price `gorm:"embedded;embeddedPrefix:row_discount_"`

	Product OrderProduct `gorm:"embedded;"`

	OrderID string

	IsGiveaway bool `gorm:"comment:是否赠品"`

	BizPayload data.JSONMap
}

func NewOrderItemFromRowDiscount(
	product OrderProduct,
	qty int32,
	taxRate apd.Decimal,
	rowDiscount price.Price,
	originalPrice price.Price,
	isGiveaway bool,
) (*OrderItem, error) {

	i := &OrderItem{
		Qty:           qty,
		RowDiscount:   rowDiscount,
		OriginalPrice: originalPrice,
		Product:       product,
		IsGiveaway:    isGiveaway,
	}

	avgDiscount, err := i.RowDiscount.ToCurrency().Div(strconv.FormatInt(int64(i.Qty), 10))
	if err != nil {
		return nil, err
	}

	pricee, err := i.OriginalPrice.ToCurrency().Sub(avgDiscount)
	if err != nil {
		return nil, err
	}

	i.Price, err = price.NewPriceFromCurrency(pricee)
	if err != nil {
		return nil, err
	}

	err = i.Cal(taxRate)

	return i, err
}

func NewOrderItemFromPriceAndOriginalPrice(
	product OrderProduct,
	qty int32,
	taxRate apd.Decimal,
	pricee price.Price,
	originalPrice price.Price,
	isGiveaway bool,
) (*OrderItem, error) {
	i := &OrderItem{
		Qty:           qty,
		Price:         pricee,
		OriginalPrice: originalPrice,
		Product:       product,
		IsGiveaway:    isGiveaway,
	}

	singleDiscount, err := i.OriginalPrice.ToCurrency().Sub(i.Price.ToCurrency())
	if err != nil {
		return nil, err
	}
	rowDiscount, err := singleDiscount.Mul(strconv.FormatInt(int64(i.Qty), 10))
	if err != nil {
		return nil, err
	}
	i.RowDiscount, err = price.NewPriceFromCurrency(rowDiscount)
	if err != nil {
		return nil, err
	}

	err = i.Cal(taxRate)

	return i, err
}

func (i *OrderItem) Cal(taxRate apd.Decimal) (err error) {

	i.Tax, i.PriceInclTax, err = calWithTax(i.Price, taxRate)

	if err != nil {
		return err
	}

	rowTotal, err := i.Price.ToCurrency().Mul(strconv.FormatInt(int64(i.Qty), 10))
	if err != nil {
		return err
	}
	i.RowTotal, err = price.NewPriceFromCurrency(rowTotal)
	if err != nil {
		return err
	}

	i.RowTotalTax, i.RowTotalInclTax, err = calWithTax(i.RowTotal, taxRate)
	if err != nil {
		return err
	}

	i.OriginalPriceTax, i.OriginalPriceInclTax, err = calWithTax(i.OriginalPrice, taxRate)
	if err != nil {
		return err
	}
	return
}

func calWithTax(p price.Price, taxRate apd.Decimal) (tax price.Price, priceInclTax price.Price, err error) {
	pp := p.ToCurrency()
	taxx, err := pp.Mul(taxRate.String())
	if err != nil {
		return
	}
	tax, err = price.NewPriceFromCurrency(taxx)
	if err != nil {
		return
	}
	priceInclTaxx, err := pp.Add(taxx)
	if err != nil {
		return
	}
	priceInclTax, err = price.NewPriceFromCurrency(priceInclTaxx)
	return
}

func (i *OrderItem) BeforeCreate(tx *gorm.DB) error {
	if len(i.ID) == 0 {
		i.ID = shortuuid.New()
	}
	return nil
}

type OrderProduct struct {
	ProductName     string
	ProductMainPic  string
	ProductID       string `gorm:"type:char(36);index:,"`
	ProductVersion  string `gorm:"type:char(36);index:,"`
	ProductType     string `gorm:"size:128;index:,"`
	ProductSkuID    string `gorm:"type:char(36);index:,"`
	ProductSkuTitle string
}

type OrderRepo interface {
	data.Repo[Order, string, *v1.ListOrderRequest]
}

func (u *Order) BeforeCreate(tx *gorm.DB) error {
	if len(u.ID) == 0 {
		u.ID = ksuid.New().String()
	}
	return nil
}

func (u *Order) FindFlowData(payway string, flowType string) []OrderFlowData {
	return lo.Filter(u.FlowData, func(item OrderFlowData, _ int) bool {
		return item.PayWay == payway && item.FlowType == flowType
	})
}

func (u *Order) RequestPay(payway string, data map[string]interface{}) {
	u.FlowData = append(u.FlowData, OrderFlowData{
		PayWay:      payway,
		FlowType:    OrderFlowTypeRequestPay,
		InitialTime: time.Now(),
		Data:        data,
		Price:       u.TotalPriceInclTax,
	})
}

func (u *Order) ChangeToPaid(payway string, paymethod string, paidPrice price.Price, data map[string]interface{}, paidTime *time.Time) {
	u.PayWay = payway
	u.PayMethod = paymethod
	u.PaidPrice = paidPrice
	u.PaidTime = paidTime
	u.Status = OrderStatusPaid
	u.FlowData = append(u.FlowData, OrderFlowData{
		PayWay:      payway,
		FlowType:    OrderFlowTypePay,
		InitialTime: time.Now(),
		Data:        data,
		Price:       paidPrice,
	})
}

func (u *Order) RequestFund(payway string, refundPrice price.Price, data map[string]interface{}) {
	u.Status = OrderStatusRefunding
	u.FlowData = append(u.FlowData, OrderFlowData{
		PayWay:      payway,
		FlowType:    OrderFlowTypeRequestRefund,
		InitialTime: time.Now(),
		Price:       refundPrice,
		Data:        data,
	})
}

func (u *Order) ChangeToRefunded(payway string, refundedPrice price.Price, data map[string]interface{}) {
	u.Status = OrderStatusRefunded
	u.FlowData = append(u.FlowData, OrderFlowData{
		PayWay:      payway,
		FlowType:    OrderFlowTypeRefund,
		InitialTime: time.Now(),
		Price:       refundedPrice,
		Data:        data,
	})

}
