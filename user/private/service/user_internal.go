package service

import (
	"context"
	klog "github.com/go-kratos/kratos/v2/log"
	dtmsrv "github.com/go-saas/kit/dtm/service"
	"github.com/go-saas/kit/dtm/utils"
	"github.com/go-saas/kit/event"
	api2 "github.com/go-saas/kit/pkg/api"
	sapi "github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/authz/authz"
	"github.com/go-saas/kit/pkg/errors"
	"github.com/go-saas/kit/pkg/uow"
	v1 "github.com/go-saas/kit/saas/api/tenant/v1"
	v12 "github.com/go-saas/kit/saas/event/v1"
	"github.com/go-saas/kit/user/api"
	pb "github.com/go-saas/kit/user/api/user/v1"
	"github.com/go-saas/kit/user/private/biz"
	"github.com/go-saas/saas"
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
	um        *biz.UserManager
	logger    *klog.Helper
}

func NewUserInternalService(
	seeder seed.Seeder,
	producer event.Producer,
	auth authz.Service,
	trust api2.TrustedContextValidator,
	dtmHelper *dtmsrv.Helper,
	um *biz.UserManager,
	l klog.Logger,
) *UserInternalService {
	return &UserInternalService{
		producer:  producer,
		auth:      auth,
		trust:     trust,
		seeder:    seeder,
		dtmHelper: dtmHelper,
		um:        um,
		logger:    klog.NewHelper(klog.With(l, "module", "user.UserInternalService")),
	}
}

func (s *UserInternalService) CreateTenant(ctx context.Context, req *v1.CreateTenantRequest) (res *emptypb.Empty, err error) {

	if err := sapi.ErrIfUntrusted(ctx, s.trust); err != nil {
		return nil, err
	}

	barrier, err := utils.BarrierFromContext(ctx)

	if err != nil {
		return nil, err
	}
	res = &emptypb.Empty{}
	err = s.dtmHelper.BarrierUow(ctx, barrier, biz.ConnName, func(ctx context.Context) error {

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
		if err := s.seeder.Seed(ctx, seed.AddTenant(req.Id), seed.WithExtra(extra)); err != nil {
			return err
		}
		e := &v12.TenantReadyEvent{
			Id:          req.Id,
			ServiceName: api.ServiceName,
		}
		ee, _ := event.NewMessageFromProto(e)
		err = s.producer.Send(ctx, ee)
		return err
	})
	return

}

// CheckUserTenant internal api for check user tenant
func (s *UserInternalService) CheckUserTenant(ctx context.Context, req *pb.CheckUserTenantRequest) (*pb.CheckUserTenantReply, error) {
	//check permission
	if err := api2.ErrIfUntrusted(ctx, s.trust); err != nil {
		return nil, err
	}

	ok, err := s.checkUserTenantInternal(ctx, req.UserId, req.TenantId)
	if err != nil {
		return nil, err
	}

	return &pb.CheckUserTenantReply{Ok: ok}, nil
}

func (s *UserInternalService) checkUserTenantInternal(ctx context.Context, userId, tenantId string) (bool, error) {
	//change to the request tenant
	ctx = saas.NewCurrentTenant(ctx, tenantId, "")
	ok, err := s.um.IsInTenant(ctx, userId, tenantId)
	if err != nil {
		return false, err
	}
	if ok {
		//user in this tenant
		return true, nil
	}
	s.logger.Debugf("user:%s not in tenant:%s", userId, tenantId)
	//super permission check
	if _, err := s.auth.Check(ctx, authz.NewEntityResource("*", "*"), authz.AnyAction); err != nil {
		//no permission
		if errors.UnRecoverableError(err) {
			//internal server error
			s.logger.Errorf("no recover error:%v", err)
			return false, err
		}
		return false, v1.ErrorTenantForbiddenLocalized(ctx, nil, nil)
	}
	return true, nil
}
