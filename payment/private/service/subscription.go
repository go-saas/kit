package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-saas/kit/payment/private/biz"
	"github.com/go-saas/kit/pkg/authn"
	"github.com/go-saas/kit/pkg/authz/authz"
	"github.com/go-saas/kit/pkg/query"
	stripe2 "github.com/go-saas/kit/pkg/stripe"
	v13 "github.com/go-saas/kit/product/api/price/v1"
	v12 "github.com/go-saas/kit/product/api/product/v1"
	v1 "github.com/go-saas/kit/user/api/user/v1"
	"github.com/go-saas/saas/data"
	"github.com/go-saas/saas/gorm"
	"github.com/samber/lo"
	"github.com/stripe/stripe-go/v76"
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
}

var _ pb.SubscriptionServiceServer = (*SubscriptionService)(nil)

func NewSubscriptionService(
	auth authz.Service,
	userInternalSrv v1.UserInternalServiceServer,
	prodInternalSrv v12.ProductInternalServiceServer,
	stripeClient *stripeclient.API,
	subsRepo biz.SubscriptionRepo,
) *SubscriptionService {
	return &SubscriptionService{
		auth:            auth,
		userInternalSrv: userInternalSrv,
		prodInternalSrv: prodInternalSrv,
		stripeClient:    stripeClient,
		subsRepo:        subsRepo,
	}
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

func (s *SubscriptionService) CreateMySubscription(ctx context.Context, req *pb.CreateMySubscriptionRequest) (*pb.Subscription, error) {
	ui, err := authn.ErrIfUnauthenticated(ctx)
	if err != nil {
		return nil, err
	}
	uid := ui.GetId()
	_, ok := lo.Find(req.Items, func(itemParams *pb.SubscriptionItemParams) bool {
		return len(itemParams.PriceId) == 0
	})
	if ok {
		return nil, errors.BadRequest("", "")
	}
	// disable tenant filter to find products and prices from different tenant
	ctx = data.NewMultiTenancyDataFilter(ctx, false)
	ctx = data.NewAutoSetTenantId(ctx, false)
	var prices []*v13.Price
	tenantId := ""
	for i, item := range req.Items {
		price, err := s.prodInternalSrv.GetInternalPrice(ctx, &v12.GetInternalPriceRequest{Id: item.PriceId})
		if err != nil {
			return nil, err
		}
		if i == 0 {
			tenantId = price.TenantId
		} else {
			if tenantId != price.TenantId {
				return nil, pb.ErrorProductAcrossTenantLocalized(ctx, nil, nil)
			}
		}
		prices = append(prices, price)
	}
	customer, err := s.userInternalSrv.FindOrCreateStripeCustomer(ctx, &v1.FindOrCreateStripeCustomerRequest{UserId: &uid})
	if err != nil {
		return nil, err
	}
	subsNewParams := &stripe.SubscriptionParams{
		Customer: stripe2.String(customer.StripeCustomerId),
		Items: lo.Map(prices, func(t *v13.Price, _ int) *stripe.SubscriptionItemsParams {
			return &stripe.SubscriptionItemsParams{Price: t.StripePriceId}
		}),
		PaymentBehavior: stripe.String("default_incomplete"),
		//TrialPeriodDays:            nil,
	}
	subsNewParams.AddExpand("latest_invoice.payment_intent")
	subs, err := s.stripeClient.Subscriptions.New(subsNewParams)
	if err != nil {
		return nil, err
	}
	localSubs := &biz.Subscription{Provider: stripe2.ProviderName, ProviderKey: subs.ID, UserId: uid}
	localSubs.TenantId = gorm.NewTenantId(tenantId)
	MapStripeSubscription2Biz(subs, localSubs)
	err = s.subsRepo.Create(ctx, localSubs)
	if err != nil {
		return nil, err
	}
	ret := &pb.Subscription{Provider: req.Provider, ProviderKey: subs.ID}
	ret.ProviderInfo = &pb.SubscriptionProviderInfo{Stripe: &stripe2.Subscription{
		Id: subs.ID,
		LatestInvoice: &stripe2.Invoice{
			Id: subs.LatestInvoice.ID,
			PaymentIntent: &stripe2.PaymentIntent{
				Id:           subs.LatestInvoice.PaymentIntent.ID,
				ClientSecret: subs.LatestInvoice.PaymentIntent.ClientSecret,
			},
		},
	}}
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
	subs, err := s.stripeClient.Subscriptions.Cancel(g.ProviderKey, nil)
	MapStripeSubscription2Biz(subs, g)
	err = s.subsRepo.Update(ctx, subs.ID, g, nil)
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
	//TODO
	return &pb.Subscription{}, nil
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

}