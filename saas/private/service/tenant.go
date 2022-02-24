package service

import (
	"context"
	"github.com/ahmetb/go-linq/v3"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authz"
	"github.com/goxiaoy/go-saas-kit/pkg/blob"
	pb "github.com/goxiaoy/go-saas-kit/saas/api/tenant/v1"
	"github.com/goxiaoy/go-saas-kit/saas/private/biz"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TenantService struct {
	pb.UnimplementedTenantServiceServer
	useCase *biz.TenantUseCase
	auth    authz.Service
	blob    blob.Factory
}

func NewTenantService(useCase *biz.TenantUseCase, auth authz.Service, blob blob.Factory) *TenantService {
	return &TenantService{useCase: useCase, auth: auth, blob: blob}
}

func (s *TenantService) CreateTenant(ctx context.Context, req *pb.CreateTenantRequest) (*pb.Tenant, error) {

	if _, err := s.auth.Check(ctx, authz.NewEntityResource("saas.tenant", "*"), authz.CreateAction); err != nil {
		return nil, err
	}

	disPlayName := req.Name
	if req.DisplayName != "" {
		disPlayName = req.DisplayName
	}
	t := &biz.Tenant{
		Name:        req.Name,
		DisplayName: disPlayName,
		Region:      req.Region,
	}
	if err := s.useCase.Create(ctx, t); err != nil {
		return nil, err
	}

	return mapBizTenantToApi(t), nil
}
func (s *TenantService) UpdateTenant(ctx context.Context, req *pb.UpdateTenantRequest) (*pb.Tenant, error) {

	if _, err := s.auth.Check(ctx, authz.NewEntityResource("saas.tenant", req.Tenant.Id), authz.UpdateAction); err != nil {
		return nil, err
	}

	t, err := s.useCase.Get(ctx, req.Tenant.Id)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, errors.NotFound("", "")
	}
	t.Name = req.Tenant.Name
	t.DisplayName = req.Tenant.DisplayName

	var tenantConn []biz.TenantConn
	linq.From(req.Tenant.Conn).SelectT(func(t *pb.TenantConnectionString) biz.TenantConn {
		return biz.TenantConn{
			Key:   t.Key,
			Value: t.Value,
		}
	}).ToSlice(&tenantConn)

	var tenantFeature []biz.TenantFeature
	linq.From(req.Tenant.Features).SelectT(func(t *pb.TenantFeature) biz.TenantFeature {
		return biz.TenantFeature{
			Key:   t.Key,
			Value: t.Value,
		}
	}).ToSlice(&tenantFeature)
	t.Conn = tenantConn
	t.Features = tenantFeature

	if err := s.useCase.Update(ctx, t, req.UpdateMask); err != nil {
		return nil, err
	}
	return mapBizTenantToApi(t), nil
}
func (s *TenantService) DeleteTenant(ctx context.Context, req *pb.DeleteTenantRequest) (*pb.DeleteTenantReply, error) {

	if _, err := s.auth.Check(ctx, authz.NewEntityResource("saas.tenant", req.Id), authz.DeleteAction); err != nil {
		return nil, err
	}

	if err := s.useCase.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &pb.DeleteTenantReply{}, nil
}

func (s *TenantService) GetTenant(ctx context.Context, req *pb.GetTenantRequest) (*pb.Tenant, error) {
	t, err := s.useCase.FindByIdOrName(ctx, req.IdOrName)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, errors.NotFound("", "")
	}

	if _, err := s.auth.Check(ctx, authz.NewEntityResource("saas.tenant", t.ID), authz.GetAction); err != nil {
		return nil, err
	}

	return mapBizTenantToApi(t), nil
}
func (s *TenantService) ListTenant(ctx context.Context, req *pb.ListTenantRequest) (*pb.ListTenantReply, error) {

	if authResult, err := s.auth.Check(ctx, authz.NewEntityResource("saas.tenant", "*"), authz.ListAction); err != nil {
		return nil, err
	} else if !authResult.Allowed {
		return nil, errors.Forbidden("", "")
	}

	ret := &pb.ListTenantReply{}

	totalCount, filterCount, err := s.useCase.Count(ctx, req.Search, req.Filter)
	ret.TotalSize = int32(totalCount)
	ret.FilterSize = int32(filterCount)
	if err != nil {
		return ret, err
	}
	items, err := s.useCase.List(ctx, req)
	if err != nil {
		return ret, err
	}
	rItems := make([]*pb.Tenant, len(items))

	linq.From(items).SelectT(func(g *biz.Tenant) *pb.Tenant { return mapBizTenantToApi(g) }).ToSlice(&rItems)
	ret.Items = rItems
	return ret, nil
}

func mapBizTenantToApi(tenant *biz.Tenant) *pb.Tenant {
	var conns []*pb.TenantConnectionString
	linq.From(tenant.Conn).SelectT(func(con biz.TenantConn) *pb.TenantConnectionString {
		return &pb.TenantConnectionString{
			Key:   con.Key,
			Value: con.Value,
		}
	}).ToSlice(&conns)

	var features []*pb.TenantFeature
	linq.From(tenant.Features).SelectT(func(con biz.TenantFeature) *pb.TenantFeature {
		return &pb.TenantFeature{
			Key:   con.Key,
			Value: con.Value,
		}
	}).ToSlice(&features)

	res := &pb.Tenant{
		Id:          tenant.ID,
		Name:        tenant.Name,
		DisplayName: tenant.DisplayName,
		Region:      tenant.Region,
		CreatedAt:   timestamppb.New(tenant.CreatedAt),
		UpdatedAt:   timestamppb.New(tenant.UpdatedAt),
		Conn:        conns,
		Features:    features,
	}
	return res
}
