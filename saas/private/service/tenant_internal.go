package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	sapi "github.com/goxiaoy/go-saas-kit/pkg/api"
	"github.com/goxiaoy/go-saas-kit/pkg/blob"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	"github.com/goxiaoy/go-saas-kit/pkg/localize"
	pb "github.com/goxiaoy/go-saas-kit/saas/api/tenant/v1"
	"github.com/goxiaoy/go-saas-kit/saas/private/biz"
)

type TenantInternalService struct {
	pb.UnimplementedTenantInternalServiceServer `wire:"-"`
	Trusted                                     sapi.TrustedContextValidator
	UseCase                                     *biz.TenantUseCase
	App                                         *conf.AppConfig
	Blob                                        blob.Factory
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
