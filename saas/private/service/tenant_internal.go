package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
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
	UseCase *biz.TenantUseCase
	App     *conf.AppConfig
	Blob    vfs.Blob
}

func NewTenantInternalService(trusted sapi.TrustedContextValidator, useCase *biz.TenantUseCase, app *conf.AppConfig, blob vfs.Blob) *TenantInternalService {
	return &TenantInternalService{Trusted: trusted, UseCase: useCase, App: app, Blob: blob}
}

func (s *TenantInternalService) GetTenant(ctx context.Context, req *pb.GetTenantRequest) (*pb.Tenant, error) {
	if ok, err := s.Trusted.Trusted(ctx); err != nil {
		return nil, err
	} else if !ok {
		//internal api call
		return nil, errors.Forbidden("", "")
	}
	t, err := s.UseCase.FindByIdOrName(ctx, req.IdOrName)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, pb.ErrorTenantNotFoundLocalized(localize.FromContext(ctx), nil, nil)
	}

	return mapBizTenantToApi(ctx, s.App, s.Blob, t), nil
}
