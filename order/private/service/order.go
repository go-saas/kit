package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-saas/kit/event"
	"github.com/go-saas/kit/order/api"
	pb "github.com/go-saas/kit/order/api/order/v1"
	"github.com/go-saas/kit/order/private/biz"
	sapi "github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/authn"
	"github.com/go-saas/kit/pkg/authz/authz"
	"github.com/go-saas/kit/pkg/price"
	"github.com/go-saas/kit/pkg/query"
	"github.com/go-saas/kit/pkg/utils"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	repo     biz.OrderRepo
	auth     authz.Service
	trust    sapi.TrustedContextValidator
	producer event.Producer
}

var _ pb.OrderServiceServer = (*OrderService)(nil)
var _ pb.OrderInternalServiceServer = (*OrderService)(nil)
var _ pb.MyOrderServiceServer = (*OrderService)(nil)

func NewOrderService(repo biz.OrderRepo, auth authz.Service, trust sapi.TrustedContextValidator, producer event.Producer) *OrderService {
	return &OrderService{repo: repo, auth: auth, trust: trust, producer: producer}
}

func (s *OrderService) ListMyOrder(ctx context.Context, req *pb.ListOrderRequest) (*pb.ListOrderReply, error) {
	userInfo, err := authn.ErrIfUnauthenticated(ctx)
	if err != nil {
		return nil, err
	}
	ret := &pb.ListOrderReply{}
	if req.Filter == nil {
		req.Filter = &pb.OrderFilter{}
	}
	req.Filter.CustomerId = &query.StringFilterOperation{Eq: &wrapperspb.StringValue{Value: userInfo.GetId()}}
	cursorRet, err := s.repo.ListCursor(ctx, req)
	if err != nil {
		return nil, err
	}
	ret.NextBeforePageToken = cursorRet.Before
	ret.NextAfterPageToken = cursorRet.After

	if err != nil {
		return ret, err
	}
	items := cursorRet.Items
	rItems := lo.Map(items, func(g *biz.Order, _ int) *pb.Order {
		b := &pb.Order{}
		MapBizOrder2Pb(ctx, g, b)
		return b
	})

	ret.Items = rItems
	return ret, nil
}

func (s *OrderService) GetMyOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error) {
	userInfo, err := authn.ErrIfUnauthenticated(ctx)
	if err != nil {
		return nil, err
	}
	g, err := s.repo.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	if g == nil || g.CustomerID != userInfo.GetId() {
		return nil, errors.NotFound("", "")
	}
	res := &pb.Order{}
	MapBizOrder2Pb(ctx, g, res)
	return res, nil
}

func (s *OrderService) RefundMyOrder(ctx context.Context, req *pb.RefundMyOrderRequest) (*pb.Order, error) {
	userInfo, err := authn.ErrIfUnauthenticated(ctx)
	if err != nil {
		return nil, err
	}
	g, err := s.repo.Get(ctx, req.GetOrderId())
	if err != nil {
		return nil, err
	}
	if g == nil || g.CustomerID != userInfo.GetId() {
		return nil, errors.NotFound("", "")
	}
	//call payment to refund
	g.RequestFund(g.PayProvider, g.TotalPriceAmount, nil)
	//TODO
	return nil, errors.BadRequest("", "")
}

func (s *OrderService) ListOrder(ctx context.Context, req *pb.ListOrderRequest) (*pb.ListOrderReply, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceOrder, "*"), authz.ReadAction); err != nil {
		return nil, err
	}
	ret := &pb.ListOrderReply{}

	totalCount, filterCount, err := s.repo.Count(ctx, req)
	ret.TotalSize = int32(totalCount)
	ret.FilterSize = int32(filterCount)

	if err != nil {
		return ret, err
	}
	items, err := s.repo.List(ctx, req)
	if err != nil {
		return ret, err
	}
	rItems := lo.Map(items, func(g *biz.Order, _ int) *pb.Order {
		b := &pb.Order{}
		MapBizOrder2Pb(ctx, g, b)
		return b
	})

	ret.Items = rItems
	return ret, nil
}
func (s *OrderService) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceOrder, req.Id), authz.ReadAction); err != nil {
		return nil, err
	}
	g, err := s.repo.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	if g == nil {
		return nil, errors.NotFound("", "")
	}
	res := &pb.Order{}
	MapBizOrder2Pb(ctx, g, res)
	return res, nil
}
func (s *OrderService) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceOrder, "*"), authz.WriteAction); err != nil {
		return nil, err
	}
	e := &biz.Order{}
	MapCreatePbOrder2Biz(req, e)
	err := s.repo.Create(ctx, e)
	if err != nil {
		return nil, err
	}
	res := &pb.Order{}
	MapBizOrder2Pb(ctx, e, res)
	return res, nil
}
func (s *OrderService) UpdateOrder(ctx context.Context, req *pb.UpdateOrderRequest) (*pb.Order, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceOrder, req.Order.Id), authz.WriteAction); err != nil {
		return nil, err
	}
	g, err := s.repo.Get(ctx, req.Order.Id)
	if err != nil {
		return nil, err
	}
	if g == nil {
		return nil, errors.NotFound("", "")
	}

	MapUpdatePbOrder2Biz(req.Order, g)
	if err := s.repo.Update(ctx, g.ID, g, nil); err != nil {
		return nil, err
	}
	res := &pb.Order{}
	MapBizOrder2Pb(ctx, g, res)
	return res, nil
}
func (s *OrderService) DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.DeleteOrderReply, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceOrder, req.Id), authz.WriteAction); err != nil {
		return nil, err
	}
	g, err := s.repo.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if g == nil {
		return nil, errors.NotFound("", "")
	}

	err = s.repo.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteOrderReply{Id: g.ID}, nil
}

func MapBizOrder2Pb(ctx context.Context, a *biz.Order, b *pb.Order) {
	b.Id = a.ID

	b.Status = a.Status
	b.CreatedAt = utils.Time2Timepb(&a.CreatedAt)

	b.TotalPrice = price.MustNewFromInt64(a.TotalPriceAmount, a.CurrencyCode).ToPricePb(ctx)
	b.TotalPriceInclTax = price.MustNewFromInt64(a.TotalPriceInclTaxAmount, a.CurrencyCode).ToPricePb(ctx)
	b.Discount = price.MustNewFromInt64(a.DiscountAmount, a.CurrencyCode).ToPricePb(ctx)
	b.OriginalPrice = price.MustNewFromInt64(a.OriginalPriceAmount, a.CurrencyCode).ToPricePb(ctx)

	b.PaidPrice = price.MustNewFromInt64(a.PaidPriceAmount, a.CurrencyCode).ToPricePb(ctx)
	b.PaidTime = utils.Time2Timepb(a.PaidTime)
	b.PayBefore = utils.Time2Timepb(a.PayBefore)
	b.PayProvider = a.PayProvider

	b.CustomerId = a.CustomerID
	b.Items = lo.Map(a.Items, func(item biz.OrderItem, _ int) *pb.OrderItem {
		return MapBizOrderItem2Pb(ctx, &item)
	})
	b.PaymentProviders = lo.Map(a.PaymentProviders, func(t biz.OrderPaymentProvider, _ int) *pb.OrderPaymentProvider {
		return MapBizOrderPaymentProvider2Pb(&t)
	})
}

func MapBizOrderItem2Pb(ctx context.Context, a *biz.OrderItem) *pb.OrderItem {
	return &pb.OrderItem{
		Id:  a.ID,
		Qty: a.Qty,

		Price:           price.MustNewFromInt64(a.PriceAmount, a.CurrencyCode).ToPricePb(ctx),
		PriceTax:        price.MustNewFromInt64(a.PriceTaxAmount, a.CurrencyCode).ToPricePb(ctx),
		PriceInclTax:    price.MustNewFromInt64(a.PriceInclTaxAmount, a.CurrencyCode).ToPricePb(ctx),
		RowTotal:        price.MustNewFromInt64(a.RowTotalAmount, a.CurrencyCode).ToPricePb(ctx),
		RowTotalTax:     price.MustNewFromInt64(a.RowTotalTaxAmount, a.CurrencyCode).ToPricePb(ctx),
		RowTotalInclTax: price.MustNewFromInt64(a.RowTotalInclTaxAmount, a.CurrencyCode).ToPricePb(ctx),
		OriginalPrice:   price.MustNewFromInt64(a.OriginalPriceAmount, a.CurrencyCode).ToPricePb(ctx),
		RowDiscount:     price.MustNewFromInt64(a.RowDiscountAmount, a.CurrencyCode).ToPricePb(ctx),
		Product: &pb.OrderProduct{
			Name:     a.Product.ProductName,
			MainPic:  a.Product.ProductMainPic,
			Id:       a.Product.ProductID,
			Version:  a.Product.ProductVersion,
			Type:     a.Product.ProductType,
			SkuId:    a.Product.ProductSkuID,
			SkuTitle: a.Product.ProductSkuTitle,
			PriceId:  a.Product.PriceID,
		},
		IsGiveaway: a.IsGiveaway,
		BizPayload: utils.Map2Structpb(a.BizPayload),
	}
}

func MapBizOrderPaymentProvider2Pb(a *biz.OrderPaymentProvider) *pb.OrderPaymentProvider {
	return &pb.OrderPaymentProvider{
		Provider:    a.Provider,
		ProviderKey: a.ProviderKey,
	}
}
func MapUpdatePbOrder2Biz(a *pb.UpdateOrder, b *biz.Order) {

}
func MapCreatePbOrder2Biz(a *pb.CreateOrderRequest, b *biz.Order) {

}
