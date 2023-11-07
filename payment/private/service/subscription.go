package service

import (
	"context"
	"github.com/go-saas/kit/event"
	v13 "github.com/go-saas/kit/payment/event/v1"
	"github.com/go-saas/kit/payment/private/biz"
	sapi "github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/authn"
	"github.com/go-saas/kit/pkg/authz/authz"
	"github.com/go-saas/kit/pkg/query"
	stripe2 "github.com/go-saas/kit/pkg/stripe"
	"github.com/go-saas/kit/pkg/utils"
	v12 "github.com/go-saas/kit/product/api/product/v1"
	v1 "github.com/go-saas/kit/user/api/user/v1"
	"github.com/go-saas/saas/data"
	"github.com/samber/lo"
	stripeclient "github.com/stripe/stripe-go/v76/client"
	"google.golang.org/protobuf/types/known/wrapperspb"

	pb "github.com/go-saas/kit/payment/api/subscription/v1"
)

type SubscriptionService struct {
	auth            authz.Service
	userInternalSrv v1.UserInternalServiceServer
	prodInternalSrv v12.ProductInternalServiceServer
	stripeClient    *stripeclient.API
	subsRepo        biz.SubscriptionRepo
	trusted         sapi.TrustedContextValidator
}

var _ pb.SubscriptionServiceServer = (*SubscriptionService)(nil)
var _ pb.SubscriptionInternalServiceServer = (*SubscriptionService)(nil)

func NewSubscriptionService(
	auth authz.Service,
	userInternalSrv v1.UserInternalServiceServer,
	prodInternalSrv v12.ProductInternalServiceServer,
	stripeClient *stripeclient.API,
	subsRepo biz.SubscriptionRepo,
	trusted sapi.TrustedContextValidator,
) *SubscriptionService {
	return &SubscriptionService{
		auth:            auth,
		userInternalSrv: userInternalSrv,
		prodInternalSrv: prodInternalSrv,
		stripeClient:    stripeClient,
		subsRepo:        subsRepo,
		trusted:         trusted,
	}
}

func (s *SubscriptionService) CreateSubscription(ctx context.Context, req *pb.CreateSubscriptionRequest) (*pb.Subscription, error) {
	//TODO
	return &pb.Subscription{}, nil
}

func (s *SubscriptionService) UpdateSubscription(ctx context.Context, req *pb.UpdateSubscriptionRequest) (*pb.Subscription, error) {
	//TODO
	return &pb.Subscription{}, nil
}

func (s *SubscriptionService) ListSubscription(ctx context.Context, req *pb.ListSubscriptionRequest) (*pb.ListSubscriptionReply, error) {
	//TODO
	return &pb.ListSubscriptionReply{}, nil
}
func (s *SubscriptionService) GetSubscription(ctx context.Context, req *pb.GetSubscriptionRequest) (*pb.Subscription, error) {
	//TODO
	return &pb.Subscription{}, nil
}

func (s *SubscriptionService) CancelSubscription(ctx context.Context, req *pb.CancelSubscriptionRequest) (*pb.Subscription, error) {
	//TODO
	return &pb.Subscription{}, nil
}

func (s *SubscriptionService) GetMySubscription(ctx context.Context, req *pb.GetMySubscriptionRequest) (*pb.Subscription, error) {
	ui, err := authn.ErrIfUnauthenticated(ctx)
	if err != nil {
		return nil, err
	}
	// disable tenant filter
	ctx = data.NewMultiTenancyDataFilter(ctx, false)
	ctx = data.NewAutoSetTenantId(ctx, false)

	g, err := s.subsRepo.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	if g == nil || g.UserId != ui.GetId() {
		return nil, pb.ErrorSubscriptionNotFoundLocalized(ctx, nil, nil)
	}
	subs, err := s.stripeClient.Subscriptions.Get(g.ProviderKey, nil)
	ret := &pb.Subscription{}
	mapBizSubscription2Pb(g, ret)
	infoSubs := &stripe2.Subscription{}
	stripe2.MapStripeSubscription(subs, infoSubs)
	ret.ProviderInfo = &pb.SubscriptionProviderInfo{Stripe: &pb.SubscriptionStripeInfo{Subscription: infoSubs}}
	return ret, nil
}

func (s *SubscriptionService) CancelMySubscription(ctx context.Context, req *pb.CancelSubscriptionRequest) (*pb.Subscription, error) {
	ui, err := authn.ErrIfUnauthenticated(ctx)
	if err != nil {
		return nil, err
	}
	// disable tenant filter
	ctx = data.NewMultiTenancyDataFilter(ctx, false)
	ctx = data.NewAutoSetTenantId(ctx, false)

	g, err := s.subsRepo.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	if g == nil || g.UserId != ui.GetId() {
		return nil, pb.ErrorSubscriptionNotFoundLocalized(ctx, nil, nil)
	}
	//TODO other providers?
	subs, err := s.stripeClient.Subscriptions.Cancel(g.ProviderKey, nil)
	MapStripeSubscription2Biz(subs, g)
	ee, _ := event.NewMessageFromProto(&v13.SubscriptionChangedEvent{Id: g.ID.String()})
	g.AppendEvent(ee)
	err = s.subsRepo.Update(ctx, g.ID.String(), g, nil)
	if err != nil {
		return nil, err
	}
	return &pb.Subscription{}, nil
}

func (s *SubscriptionService) UpdateMySubscription(ctx context.Context, req *pb.UpdateMySubscriptionRequest) (*pb.Subscription, error) {
	ui, err := authn.ErrIfUnauthenticated(ctx)
	if err != nil {
		return nil, err
	}
	// disable tenant filter to find subs
	ctx = data.NewMultiTenancyDataFilter(ctx, false)
	ctx = data.NewAutoSetTenantId(ctx, false)

	g, err := s.subsRepo.Get(ctx, req.Subscription.GetId())
	if err != nil {
		return nil, err
	}
	if g == nil || g.UserId != ui.GetId() {
		return nil, pb.ErrorSubscriptionNotFoundLocalized(ctx, nil, nil)
	}
	//TODO update subscription
	ret := &pb.Subscription{}
	mapBizSubscription2Pb(g, ret)

	return ret, nil
}

func (s *SubscriptionService) ListMySubscription(ctx context.Context, req *pb.ListMySubscriptionRequest) (*pb.ListMySubscriptionReply, error) {
	ui, err := authn.ErrIfUnauthenticated(ctx)
	if err != nil {
		return nil, err
	}
	if req.Filter == nil {
		req.Filter = &pb.SubscriptionFilter{}
	}
	req.Filter.And = append(req.Filter.And, &pb.SubscriptionFilter{UserId: &query.StringFilterOperation{Eq: wrapperspb.String(ui.GetId())}})

	ctx = data.NewMultiTenancyDataFilter(ctx, false)

	ret := &pb.ListMySubscriptionReply{}

	totalCount, filterCount, err := s.subsRepo.Count(ctx, req)
	ret.TotalSize = int32(totalCount)
	ret.FilterSize = int32(filterCount)
	if err != nil {
		return ret, err
	}
	items, err := s.subsRepo.List(ctx, req)
	if err != nil {
		return ret, err
	}
	rItems := lo.Map(items, func(g *biz.Subscription, _ int) *pb.Subscription {
		b := &pb.Subscription{}
		mapBizSubscription2Pb(g, b)
		return b
	})
	ret.Items = rItems
	return ret, nil

}

func mapBizSubscription2Pb(a *biz.Subscription, b *pb.Subscription) {
	b.Id = a.ID.String()
	b.Provider = a.Provider
	b.ProviderKey = a.ProviderKey
	b.UserId = a.UserId
	b.Status = string(a.Status)
	b.Items = lo.Map(a.Items, func(t biz.SubscriptionItem, i int) *pb.SubscriptionItem {
		ret := &pb.SubscriptionItem{}
		mapBizSubscriptionItem2Pb(&t, ret)
		return ret
	})
}

func mapBizSubscriptionItem2Pb(a *biz.SubscriptionItem, b *pb.SubscriptionItem) {
	b.Id = a.ID.String()
	b.PriceId = a.PriceID
	b.ProductId = a.ProductID
	b.PriceOwnerId = a.PriceOwnerID
	b.PriceOwnerType = a.PriceOwnerType
	b.Quantity = a.Quantity
	b.BizPayload = utils.Map2Structpb(a.BizPayload)
}
