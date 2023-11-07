package service

import (
	"context"
	"github.com/cockroachdb/apd/v3"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-saas/kit/event"
	pb "github.com/go-saas/kit/order/api/order/v1"
	v1 "github.com/go-saas/kit/order/event/v1"
	"github.com/go-saas/kit/order/private/biz"
	"github.com/go-saas/kit/pkg/utils"
	"github.com/go-saas/lbs"
	"time"
)

func (s *OrderService) CreateInternalOrder(ctx context.Context, req *pb.CreateInternalOrderRequest) (*pb.Order, error) {
	if ok, _ := s.trust.Trusted(ctx); !ok {
		return nil, errors.Forbidden("", "")
	}
	taxRate, _, _ := apd.NewFromString("0")
	var orderItems []biz.OrderItem
	for _, item := range req.Items {
		orderItem, err := biz.NewOrderItemFromPriceAndOriginalPrice(
			req.CurrencyCode,
			biz.OrderProduct{
				ProductName:     item.Product.Name,
				ProductMainPic:  item.Product.MainPic,
				ProductID:       item.Product.Id,
				ProductVersion:  item.Product.Version,
				ProductType:     item.Product.Type,
				ProductSkuID:    item.Product.SkuId,
				ProductSkuTitle: item.Product.SkuTitle,
			},
			item.Qty,
			*taxRate,
			item.PriceAmount,
			item.OriginalPriceAmount,
			item.IsGiveaway,
			utils.Structpb2Map(item.BizPayload),
		)
		if err != nil {
			return nil, err
		}
		orderItems = append(orderItems, *orderItem)
	}
	e, err := biz.NewOrder(req.CurrencyCode, *taxRate, orderItems)
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
	if req.Id == nil && (req.Provider == nil || req.ProviderKey == nil) {
		return nil, errors.BadRequest("", "provider id or (provider and providerKey)")
	}
	var g *biz.Order
	var err error
	if req.Id != nil {
		g, err = s.repo.Get(ctx, req.GetId())
	} else {
		g, err = s.repo.FindByPaymentProvider(ctx, req.GetProvider(), req.GetProviderKey())
	}
	g, err = s.repo.Get(ctx, req.GetId())
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

	g.ChangeToPaid(req.PayProvider, req.PayMethod, req.PaidPriceAmount, utils.Timepb2Time(req.PaidTime))
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

	g.ChangeToRefunded(req.PayProvider, req.RefundPriceAmount, utils.Structpb2Map(req.PayExtra))
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

func (s *OrderService) UpdateInternalOrderPaymentProvider(ctx context.Context, req *pb.UpdateInternalOrderPaymentProviderRequest) (*pb.Order, error) {
	if ok, _ := s.trust.Trusted(ctx); !ok {
		return nil, errors.Forbidden("", "")
	}
	g, err := s.repo.Get(ctx, req.OrderId)
	if err != nil {
		return nil, err
	}
	if g == nil {
		return nil, errors.NotFound("", "")
	}
	paymentProvider := &biz.OrderPaymentProvider{
		OrderID:     g.ID,
		Provider:    req.Provider,
		ProviderKey: req.ProviderKey,
	}
	err = s.repo.UpsertPaymentProvider(ctx, g, paymentProvider)
	if err != nil {
		return nil, err
	}
	orderPb := &pb.Order{}
	MapBizOrder2Pb(ctx, g, orderPb)
	return orderPb, nil
}
