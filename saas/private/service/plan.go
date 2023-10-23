package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-saas/kit/pkg/authz/authz"
	"github.com/go-saas/kit/pkg/query"
	"github.com/go-saas/kit/saas/api"
	"github.com/go-saas/kit/saas/private/biz"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"strings"

	pb "github.com/go-saas/kit/saas/api/plan/v1"
)

type PlanService struct {
	pb.UnimplementedPlanServiceServer
	auth   authz.Service
	repo   biz.PlanRepo
	logger *klog.Helper
}

func NewPlanService(auth authz.Service, repo biz.PlanRepo, logger klog.Logger) *PlanService {
	return &PlanService{auth: auth, repo: repo, logger: klog.NewHelper(klog.With(logger, "module", "PlanService"))}
}

func (s *PlanService) ListPlan(ctx context.Context, req *pb.ListPlanRequest) (*pb.ListPlanReply, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourcePlan, "*"), authz.ReadAction); err != nil {
		return nil, err
	}
	ret := &pb.ListPlanReply{}

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
	rItems := lo.Map(items, func(g *biz.Plan, _ int) *pb.Plan {
		b := &pb.Plan{}
		MapBizPlan2Pb(g, b)
		return b
	})

	ret.Items = rItems
	return ret, nil
}

func (s *PlanService) GetPlan(ctx context.Context, req *pb.GetPlanRequest) (*pb.Plan, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourcePlan, req.Key), authz.ReadAction); err != nil {
		return nil, err
	}
	g, err := s.repo.Get(ctx, req.GetKey())
	if err != nil {
		return nil, err
	}
	if g == nil {
		return nil, errors.NotFound("", "")
	}
	res := &pb.Plan{}
	MapBizPlan2Pb(g, res)
	return res, nil
}
func (s *PlanService) CreatePlan(ctx context.Context, req *pb.CreatePlanRequest) (*pb.Plan, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourcePlan, "*"), authz.CreateAction); err != nil {
		return nil, err
	}

	//check duplicate name
	if dbP, err := s.repo.Get(ctx, normalizeName(req.Key)); err != nil {
		return nil, err
	} else if dbP != nil {
		return nil, pb.ErrorDuplicatePlanKeyLocalized(ctx, nil, nil)
	}
	e := &biz.Plan{}
	MapCreatePbPlan2Biz(req, e)
	err := s.repo.Create(ctx, e)
	if err != nil {
		return nil, err
	}
	res := &pb.Plan{}
	MapBizPlan2Pb(e, res)
	return res, nil
}
func (s *PlanService) UpdatePlan(ctx context.Context, req *pb.UpdatePlanRequest) (*pb.Plan, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourcePlan, req.Plan.Key), authz.UpdateAction); err != nil {
		return nil, err
	}

	//check duplicate name
	if dbP, err := s.repo.Get(ctx, normalizeName(req.Plan.Key)); err != nil {
		return nil, err
	} else if dbP != nil && dbP.Key != req.Plan.Key {
		return nil, pb.ErrorDuplicatePlanKeyLocalized(ctx, nil, nil)
	}

	g, err := s.repo.Get(ctx, req.Plan.Key)
	if err != nil {
		return nil, err
	}
	if g == nil {
		return nil, errors.NotFound("", "")
	}
	//copy plan
	MapUpdatePbPlan2Biz(req.Plan, g)
	if err := s.repo.Update(ctx, req.Plan.Key, g, nil); err != nil {
		return nil, err
	}
	res := &pb.Plan{}
	MapBizPlan2Pb(g, res)
	return res, nil
}

func (s *PlanService) DeletePlan(ctx context.Context, req *pb.DeletePlanRequest) (*pb.DeletePlanReply, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourcePlan, req.Key), authz.DeleteAction); err != nil {
		return nil, err
	}
	g, err := s.repo.Get(ctx, req.Key)
	if err != nil {
		return nil, err
	}
	if g == nil {
		return nil, errors.NotFound("", "")
	}

	err = s.repo.Delete(ctx, req.Key)
	if err != nil {
		return nil, err
	}
	return &pb.DeletePlanReply{}, nil
}

func (s *PlanService) GetAvailablePlans(ctx context.Context, req *pb.GetAvailablePlansRequest) (*pb.GetAvailablePlansReply, error) {
	items, err := s.repo.List(ctx, &pb.ListPlanRequest{
		PageOffset: 0,
		PageSize:   -1,
		Filter:     &pb.PlanFilter{Active: &query.BooleanFilterOperators{Eq: &wrapperspb.BoolValue{Value: true}}},
	})
	if err != nil {
		return nil, err
	}
	var retItems = lo.Map(items, func(a *biz.Plan, _ int) *pb.Plan {
		ret := &pb.Plan{}
		MapBizPlan2Pb(a, ret)
		return ret
	})
	return &pb.GetAvailablePlansReply{Items: retItems}, nil
}

func MapBizPlan2Pb(a *biz.Plan, b *pb.Plan) {
	b.Key = a.Key
	b.DisplayName = a.DisplayName
	b.Active = a.Active
}

func MapUpdatePbPlan2Biz(a *pb.UpdatePlan, b *biz.Plan) {
	b.DisplayName = a.DisplayName
	b.Active = a.Active
}

func MapCreatePbPlan2Biz(a *pb.CreatePlanRequest, b *biz.Plan) {
	b.Key = a.Key
	b.DisplayName = a.DisplayName
}

func normalizeName(name string) string {
	return strings.ToLower(name)
}
