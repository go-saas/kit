package service

import (
	"context"
	pb "github.com/go-saas/kit/payment/api/subscription/v1"
	sapi "github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/saas/data"
)

func (s *SubscriptionService) GetInternalSubscription(ctx context.Context, req *pb.GetInternalSubscriptionRequest) (*pb.Subscription, error) {
	if err := sapi.ErrIfUntrusted(ctx, s.trusted); err != nil {
		return nil, err
	}
	// disable tenant filter
	ctx = data.NewMultiTenancyDataFilter(ctx, false)
	g, err := s.subsRepo.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	if g == nil {
		return nil, pb.ErrorSubscriptionNotFoundLocalized(ctx, nil, nil)
	}
	ret := &pb.Subscription{}
	mapBizSubscription2Pb(g, ret)
	return ret, nil
}
