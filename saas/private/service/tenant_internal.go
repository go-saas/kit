package service

import (
	"context"
	"github.com/go-saas/kit/dtm/utils"
	sapi "github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/conf"
	"github.com/go-saas/kit/pkg/localize"
	pb "github.com/go-saas/kit/saas/api/tenant/v1"
	"github.com/go-saas/kit/saas/private/biz"
	"github.com/go-saas/uow"
	"github.com/google/uuid"
	"github.com/goxiaoy/vfs"
)

type TenantInternalService struct {
	pb.UnimplementedTenantInternalServiceServer
	Trusted sapi.TrustedContextValidator
	useCase *biz.TenantUseCase
	app     *conf.AppConfig
	blob    vfs.Blob
	uowMgr  uow.Manager
}

func NewTenantInternalService(
	trusted sapi.TrustedContextValidator,
	useCase *biz.TenantUseCase,
	app *conf.AppConfig,
	blob vfs.Blob,
	uowMgr uow.Manager,
) *TenantInternalService {
	return &TenantInternalService{Trusted: trusted, useCase: useCase, app: app, blob: blob, uowMgr: uowMgr}
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

func (s *TenantInternalService) CreateTenant(ctx context.Context, req *pb.CreateTenantRequest) (res *pb.Tenant, err error) {
	if err := sapi.ErrIfUntrusted(ctx, s.Trusted); err != nil {
		return nil, err
	}
	// TODO XA Transaction
	_, err = utils.BarrierFromContext(ctx)
	if err != nil {
		return
	}

	if len(req.DisplayName) == 0 {
		req.DisplayName = req.Name
	}
	t := &biz.Tenant{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Region:      req.Region,
		Logo:        req.Logo,
		SeparateDb:  req.SeparateDb,
	}
	if len(req.Id) > 0 {
		t.UIDBase.ID = uuid.MustParse(req.Id)
	}

	if err := s.useCase.Create(ctx, t); err != nil {
		return nil, err
	}
	res = mapBizTenantToApi(ctx, s.app, s.blob, t)
	return
}
