package service

import (
	"context"
	"github.com/cockroachdb/apd/v3"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-saas/kit/event"
	"github.com/go-saas/kit/order/api"
	pb "github.com/go-saas/kit/order/api/order/v1"
	v1 "github.com/go-saas/kit/order/event/v1"
	"github.com/go-saas/kit/order/private/biz"
	sapi "github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/authn"
	"github.com/go-saas/kit/pkg/authz/authz"
	"github.com/go-saas/kit/pkg/price"
	"github.com/go-saas/kit/pkg/query"
	"github.com/go-saas/kit/pkg/utils"
	"github.com/go-saas/lbs"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"time"
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
var _ pb.OrderAppServiceServer = (*OrderService)(nil)

func NewOrderService(repo biz.OrderRepo, auth authz.Service, trust sapi.TrustedContextValidator, producer event.Producer) *OrderService {
	return &OrderService{repo: repo, auth: auth, trust: trust, producer: producer}
}

func (s *OrderService) ListAppOrder(ctx context.Context, req *pb.ListOrderRequest) (*pb.ListOrderReply, error) {
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

func (s *OrderService) GetAppOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error) {
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

func (s *OrderService) RefundAppOrder(ctx context.Context, req *pb.RefundAppOrderRequest) (*pb.Order, error) {
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
	g.RequestFund(g.PayWay, g.TotalPriceInclTax)
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

func (s *OrderService) CreateInternalOrder(ctx context.Context, req *pb.CreateInternalOrderRequest) (*pb.Order, error) {
	if ok, _ := s.trust.Trusted(ctx); !ok {
		return nil, errors.Forbidden("", "")
	}
	taxRate, _, _ := apd.NewFromString("0")
	var orderItems []biz.OrderItem
	for _, item := range req.Items {
		pricee, err := price.NewPriceFromPb(item.Price)
		if err != nil {
			return nil, err
		}
		originalPrice, err := price.NewPriceFromPb(item.OriginalPrice)
		if err != nil {
			return nil, err
		}

		orderItem, err := biz.NewOrderItemFromPriceAndOriginalPrice(biz.OrderProduct{
			ProductName:     item.Product.Name,
			ProductMainPic:  item.Product.MainPic,
			ProductID:       item.Product.Id,
			ProductVersion:  item.Product.Version,
			ProductType:     item.Product.Type,
			ProductSkuID:    item.Product.SkuId,
			ProductSkuTitle: item.Product.SkuTitle,
		}, item.Qty, *taxRate, pricee, originalPrice, item.IsGiveaway)
		if err != nil {
			return nil, err
		}
		orderItems = append(orderItems, *orderItem)
	}
	e, err := biz.NewOrder(*taxRate, orderItems)
	if err != nil {
		return nil, err
	}
	e.CustomerID = req.CustomerId
	e.Extra = utils.Structpb2Map(req.Extra)

	if req.BillingAddr != nil {
		billingAddr, _ := lbs.NewAddressEntityFromPb(req.BillingAddr)
		e.BillingAddr = *billingAddr
	}
	if req.ShippingAddr != nil {
		shippingAddr, _ := lbs.NewAddressEntityFromPb(req.ShippingAddr)
		e.ShippingAddr = *shippingAddr
	}

	if req.PayBefore != nil {
		t := time.Now().Add(req.PayBefore.AsDuration())
		e.PayBefore = &t
	}

	err = s.repo.Create(ctx, e)
	if err != nil {
		return nil, err
	}
	res := &pb.Order{}
	MapBizOrder2Pb(ctx, e, res)
	return res, nil
}

func (s *OrderService) GetInternalOrder(ctx context.Context, req *pb.GetInternalOrderRequest) (*pb.Order, error) {
	if ok, _ := s.trust.Trusted(ctx); !ok {
		return nil, errors.Forbidden("", "")
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

func (s *OrderService) InternalOrderPaySuccess(ctx context.Context, req *pb.InternalOrderPaySuccessRequest) (*pb.InternalOrderPaySuccessReply, error) {
	if ok, _ := s.trust.Trusted(ctx); !ok {
		return nil, errors.Forbidden("", "")
	}
	g, err := s.repo.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	if g == nil {
		return nil, errors.NotFound("", "")
	}
	p, err := price.NewPriceFromPb(req.PaidPrice)
	if err != nil {
		return nil, err
	}
	g.ChangeToPaid(req.PayWay, req.PayMethod, p, utils.Structpb2Map(req.PayExtra), utils.Timepb2Time(req.PaidTime))
	if err := s.repo.Update(ctx, g.ID, g, nil); err != nil {
		return nil, err
	}
	//publish event
	orderPb := &pb.Order{}
	MapBizOrder2Pb(ctx, g, orderPb)
	msg, _ := event.NewMessageFromProto(&v1.OrderPaySuccessEvent{Order: orderPb})
	err = s.producer.Send(ctx, msg)
	if err != nil {
		return nil, err
	}
	return &pb.InternalOrderPaySuccessReply{}, err
}

func (s *OrderService) InternalOrderRefunded(ctx context.Context, req *pb.InternalOrderRefundedRequest) (*pb.InternalOrderRefundedReply, error) {
	if ok, _ := s.trust.Trusted(ctx); !ok {
		return nil, errors.Forbidden("", "")
	}
	g, err := s.repo.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	if g == nil {
		return nil, errors.NotFound("", "")
	}
	p, err := price.NewPriceFromPb(req.RefundPrice)
	if err != nil {
		return nil, err
	}
	g.ChangeToRefunded(req.PayWay, p, utils.Structpb2Map(req.PayExtra))
	if err := s.repo.Update(ctx, g.ID, g, nil); err != nil {
		return nil, err
	}
	//publish event
	orderPb := &pb.Order{}
	MapBizOrder2Pb(ctx, g, orderPb)
	msg, _ := event.NewMessageFromProto(&v1.OrderRefundSuccessEvent{Order: orderPb})
	err = s.producer.Send(ctx, msg)
	if err != nil {
		return nil, err
	}
	return &pb.InternalOrderRefundedReply{}, err
}

func MapBizOrder2Pb(ctx context.Context, a *biz.Order, b *pb.Order) {
	b.Id = a.ID

	b.Status = a.Status
	b.CreatedAt = utils.Time2Timepb(&a.CreatedAt)

	b.TotalPrice = a.TotalPrice.ToPricePb(ctx)
	b.TotalPriceInclTax = a.TotalPriceInclTax.ToPricePb(ctx)
	b.Discount = a.Discount.ToPricePb(ctx)
	b.OriginalPrice = a.OriginalPrice.ToPricePb(ctx)
	b.OriginalPriceInclTax = a.OriginalPriceInclTax.ToPricePb(ctx)
	b.PaidPrice = a.PaidPrice.ToPricePb(ctx)
	b.PaidTime = utils.Time2Timepb(a.PaidTime)
	b.PayBefore = utils.Time2Timepb(a.PayBefore)
	b.PayWay = a.PayWay

	b.CustomerId = a.CustomerID
	b.Items = lo.Map(a.Items, func(item biz.OrderItem, _ int) *pb.OrderItem {
		return MapBizOrderItem2Pb(ctx, &item)
	})
}

func MapBizOrderItem2Pb(ctx context.Context, a *biz.OrderItem) *pb.OrderItem {
	return &pb.OrderItem{
		Id:              a.ID,
		Qty:             a.Qty,
		Price:           a.Price.ToPricePb(ctx),
		Tax:             a.Tax.ToPricePb(ctx),
		PriceInclTax:    a.PriceInclTax.ToPricePb(ctx),
		RowTotal:        a.RowTotal.ToPricePb(ctx),
		RowTotalTax:     a.RowTotalTax.ToPricePb(ctx),
		RowTotalInclTax: a.RowTotalInclTax.ToPricePb(ctx),
		OriginalPrice:   a.OriginalPrice.ToPricePb(ctx),
		RowDiscount:     a.RowDiscount.ToPricePb(ctx),
		Product: &pb.OrderProduct{
			Name:     a.Product.ProductName,
			MainPic:  a.Product.ProductMainPic,
			Id:       a.Product.ProductID,
			Version:  a.Product.ProductVersion,
			Type:     a.Product.ProductType,
			SkuId:    a.Product.ProductSkuID,
			SkuTitle: a.Product.ProductSkuTitle,
		},
		IsGiveaway: a.IsGiveaway,
	}
}

func MapUpdatePbOrder2Biz(a *pb.UpdateOrder, b *biz.Order) {

}
func MapCreatePbOrder2Biz(a *pb.CreateOrderRequest, b *biz.Order) {

}
