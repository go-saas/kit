package service

import (
	"context"
	sapi "github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/conf"
	pb "github.com/go-saas/kit/saas/api/tenant/v1"
	"github.com/go-saas/kit/saas/private/biz"
	"github.com/goxiaoy/vfs"
	"github.com/samber/lo"
)

type TenantInternalService struct {
	Trusted sapi.TrustedContextValidator
	useCase *biz.TenantUseCase
	app     *conf.AppConfig
	blob    vfs.Blob
}

var _ pb.TenantInternalServiceServer = (*TenantInternalService)(nil)

func NewTenantInternalService(
	trusted sapi.TrustedContextValidator,
	useCase *biz.TenantUseCase,
	app *conf.AppConfig,
	blob vfs.Blob,
) *TenantInternalService {
	return &TenantInternalService{Trusted: trusted, useCase: useCase, app: app, blob: blob}
}

func (s *TenantInternalService) GetTenant(ctx context.Context, req *pb.GetTenantRequest) (*pb.Tenant, error) {
	if err := sapi.ErrIfUntrusted(ctx, s.Trusted); err != nil {
		return nil, err
	}
	t, err := s.useCase.FindByIdOrName(ctx, req.IdOrName)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, pb.ErrorTenantNotFoundLocalized(ctx, nil, nil)
	}

	return mapBizTenantToApi(ctx, s.app, s.blob, t), nil
}

func (s *TenantInternalService) ListTenant(ctx context.Context, req *pb.ListTenantRequest) (*pb.ListTenantReply, error) {
	if err := sapi.ErrIfUntrusted(ctx, s.Trusted); err != nil {
		return nil, err
	}
	ret := &pb.ListTenantReply{}
	totalCount, filterCount, err := s.useCase.Count(ctx, req)
	ret.TotalSize = int32(totalCount)
	ret.FilterSize = int32(filterCount)
	if err != nil {
		return ret, err
	}
	items, err := s.useCase.List(ctx, req)
	if err != nil {
		return ret, err
	}
	rItems := lo.Map(items, func(g *biz.Tenant, _ int) *pb.Tenant { return mapBizTenantToApi(ctx, s.app, s.blob, g) })
	ret.Items = rItems
	return ret, nil
}
