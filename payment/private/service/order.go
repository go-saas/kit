package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	pb "github.com/goxiaoy/go-saas-kit/payment/api/order/v1"
	"github.com/goxiaoy/go-saas-kit/payment/private/biz"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authz"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrderServiceService struct {
	pb.UnimplementedOrderServiceServer
	repo biz.OrderRepo
	auth authz.Service
}

func NewOrderServiceService(repo biz.OrderRepo, auth authz.Service) *OrderServiceService {
	return &OrderServiceService{repo: repo, auth: auth}
}

func (s *OrderServiceService) ListOrder(ctx context.Context, req *pb.ListOrderRequest) (*pb.ListOrderReply, error) {
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
		MapBizOrder2Pb(g, b)
		return b
	})

	ret.Items = rItems
	return ret, nil
}
func (s *OrderServiceService) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error) {
	g, err := s.repo.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	if g == nil {
		return nil, errors.NotFound("", "")
	}
	res := &pb.Order{}
	MapBizOrder2Pb(g, res)
	return res, nil
}
func (s *OrderServiceService) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	e := &biz.Order{}
	MapCreatePbOrder2Biz(req, e)
	err := s.repo.Create(ctx, e)
	if err != nil {
		return nil, err
	}
	res := &pb.Order{}
	MapBizOrder2Pb(e, res)
	return res, nil
}
func (s *OrderServiceService) UpdateOrder(ctx context.Context, req *pb.UpdateOrderRequest) (*pb.Order, error) {
	g, err := s.repo.Get(ctx, req.Order.Id)
	if err != nil {
		return nil, err
	}
	if g == nil {
		return nil, errors.NotFound("", "")
	}

	MapUpdatePbOrder2Biz(req.Order, g)
	if err := s.repo.Update(ctx, g.ID.String(), g, nil); err != nil {
		return nil, err
	}
	res := &pb.Order{}
	MapBizOrder2Pb(g, res)
	return res, nil
}
func (s *OrderServiceService) DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.DeleteOrderReply, error) {
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
	return &pb.DeleteOrderReply{Id: g.ID.String(), Name: g.Name}, nil
}
func MapBizOrder2Pb(a *biz.Order, b *pb.Order) {
	b.Id = a.ID.String()
	b.Name = a.Name
	b.CreatedAt = timestamppb.New(a.CreatedAt)
}

func MapUpdatePbOrder2Biz(a *pb.UpdateOrder, b *biz.Order) {
	b.Name = a.Name
}
func MapCreatePbOrder2Biz(a *pb.CreateOrderRequest, b *biz.Order) {
	b.Name = a.Name
}
