package service

import (
	"context"
	sapi "github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/conf"
	"github.com/go-saas/kit/pkg/localize"
	pb "github.com/go-saas/kit/saas/api/tenant/v1"
	"github.com/go-saas/kit/saas/private/biz"
	"github.com/goxiaoy/vfs"
)

type TenantInternalService struct {
	pb.UnimplementedTenantInternalServiceServer
	Trusted sapi.TrustedContextValidator
	useCase *biz.TenantUseCase
	app     *conf.AppConfig
	blob    vfs.Blob
}

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
		return nil, pb.ErrorTenantNotFoundLocalized(localize.FromContext(ctx), nil, nil)
	}

	return mapBizTenantToApi(ctx, s.app, s.blob, t), nil
}
