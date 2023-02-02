package service

import (
	"context"
	dtmsrv "github.com/go-saas/kit/dtm/service"
	"github.com/go-saas/kit/dtm/utils"
	"github.com/go-saas/kit/event"
	api2 "github.com/go-saas/kit/pkg/api"
	sapi "github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/authz/authz"
	"github.com/go-saas/kit/pkg/uow"
	v1 "github.com/go-saas/kit/saas/api/tenant/v1"
	v12 "github.com/go-saas/kit/saas/event/v1"
	"github.com/go-saas/kit/user/api"
	pb "github.com/go-saas/kit/user/api/user/v1"
	"github.com/go-saas/kit/user/private/biz"
	"github.com/go-saas/saas/seed"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserInternalService struct {
	pb.UnimplementedUserInternalServiceServer

	producer  event.Producer
	auth      authz.Service
	trust     api2.TrustedContextValidator
	seeder    seed.Seeder
	dtmHelper *dtmsrv.Helper
}

func NewUserInternalService(
	seeder seed.Seeder,
	producer event.Producer,
	auth authz.Service,
	trust api2.TrustedContextValidator,
	dtmHelper *dtmsrv.Helper,
) *UserInternalService {
	return &UserInternalService{
		producer:  producer,
		auth:      auth,
		trust:     trust,
		seeder:    seeder,
		dtmHelper: dtmHelper,
	}
}

func (u *UserInternalService) CreateTenant(ctx context.Context, req *v1.CreateTenantRequest) (res *emptypb.Empty, err error) {

	if err := sapi.ErrIfUntrusted(ctx, u.trust); err != nil {
		return nil, err
	}

	barrier, err := utils.BarrierFromContext(ctx)

	//return nil, status.New(codes.Aborted, "").Err()
	if err != nil {
		return nil, err
	}
	res = &emptypb.Empty{}
	err = u.dtmHelper.BarrierUow(ctx, barrier, biz.ConnName, func(ctx context.Context) error {

		ctx = uow.SkipUow(ctx)

		extra := map[string]interface{}{}
		if req.AdminEmail != nil {
			extra[biz.AdminEmailKey] = req.AdminEmail.Value
		}
		if req.AdminUsername != nil {
			extra[biz.AdminUsernameKey] = req.AdminUsername.Value
		}
		if req.AdminPassword != nil {
			extra[biz.AdminPasswordKey] = req.AdminPassword.Value
		}
		if req.AdminUserId != nil {
			extra[biz.AdminUserId] = req.AdminUserId.Value
		}
		if err := u.seeder.Seed(ctx, seed.AddTenant(req.Id), seed.WithExtra(extra)); err != nil {
			return err
		}
		e := &v12.TenantReadyEvent{
			Id:          req.Id,
			ServiceName: api.ServiceName,
		}
		ee, _ := event.NewMessageFromProto(e)
		err = u.producer.Send(ctx, ee)
		return err
	})
	return

}
