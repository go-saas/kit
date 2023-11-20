package service

import (
	"context"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	klog "github.com/go-kratos/kratos/v2/log"
	dtmsrv "github.com/go-saas/kit/dtm/service"
	"github.com/go-saas/kit/dtm/utils"
	"github.com/go-saas/kit/event"
	kapi "github.com/go-saas/kit/pkg/api"
	sapi "github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/authz/authz"
	"github.com/go-saas/kit/pkg/errors"
	"github.com/go-saas/kit/pkg/idp"
	"github.com/go-saas/kit/pkg/uow"
	v1 "github.com/go-saas/kit/saas/api/tenant/v1"
	v12 "github.com/go-saas/kit/saas/event/v1"
	"github.com/go-saas/kit/user/api"
	pb "github.com/go-saas/kit/user/api/user/v1"
	"github.com/go-saas/kit/user/private/biz"
	"github.com/go-saas/saas"
	"github.com/go-saas/saas/seed"
	"github.com/goxiaoy/vfs"
	"github.com/samber/lo"
	"github.com/stripe/stripe-go/v76"
	stripeclient "github.com/stripe/stripe-go/v76/client"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserInternalService struct {
	producer     event.Producer
	auth         authz.Service
	trust        kapi.TrustedContextValidator
	seeder       seed.Seeder
	dtmHelper    *dtmsrv.Helper
	blob         vfs.Blob
	um           *biz.UserManager
	logger       *klog.Helper
	stripeClient *stripeclient.API
}

var _ pb.UserInternalServiceServer = (*UserInternalService)(nil)

func NewUserInternalService(
	seeder seed.Seeder,
	producer event.Producer,
	auth authz.Service,
	trust kapi.TrustedContextValidator,
	dtmHelper *dtmsrv.Helper,
	blob vfs.Blob,
	um *biz.UserManager,
	l klog.Logger,
	stripeClient *stripeclient.API,
) *UserInternalService {
	return &UserInternalService{
		producer:     producer,
		auth:         auth,
		trust:        trust,
		seeder:       seeder,
		dtmHelper:    dtmHelper,
		blob:         blob,
		um:           um,
		stripeClient: stripeClient,
		logger:       klog.NewHelper(klog.With(l, "module", "user.UserInternalService")),
	}
}

func (s *UserInternalService) CreateTenant(ctx context.Context, req *pb.UserInternalCreateTenantRequest) (res *emptypb.Empty, err error) {

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
			extra[biz.AdminEmailKey] = *req.AdminEmail
		}
		if req.AdminUsername != nil {
			extra[biz.AdminUsernameKey] = *req.AdminUsername
		}
		if req.AdminPassword != nil {
			extra[biz.AdminPasswordKey] = *req.AdminPassword
		}
		if req.AdminUserId != nil {
			extra[biz.AdminUserId] = *req.AdminUserId
		}
		if err := s.seeder.Seed(ctx, seed.AddTenant(req.TenantId), seed.WithExtra(extra)); err != nil {
			return err
		}
		e := &v12.TenantReadyEvent{
			Id:          req.TenantId,
			ServiceName: api.ServiceName,
		}
		ee, _ := event.NewMessageFromProto(e)
		err = s.producer.Send(ctx, ee)
		return err
	})
	return

}

func (s *UserInternalService) FindOrCreateStripeCustomer(ctx context.Context, req *pb.FindOrCreateStripeCustomerRequest) (*pb.FindOrCreateStripeCustomerReply, error) {
	if err := sapi.ErrIfUntrusted(ctx, s.trust); err != nil {
		return nil, err
	}
	if req.UserId == nil && req.StripeCustomerId == nil {
		return nil, kerrors.BadRequest("", "")
	}
	ret := &pb.FindOrCreateStripeCustomerReply{}
	if req.UserId != nil {
		user, err := s.um.FindByID(ctx, *req.UserId)
		if err != nil {
			return nil, err
		}
		if user == nil {
			return nil, pb.ErrorUserNotFoundLocalized(ctx, nil, nil)
		}
		//get login providers
		logins, err := s.um.ListLogin(ctx, user)
		if err != nil {
			return nil, err
		}
		stripeLogin, ok := lo.Find(logins, func(login *biz.UserLogin) bool {
			return login.LoginProvider == idp.StripeLoginProvider
		})
		if !ok {
			params := &stripe.CustomerParams{
				Name: user.Name,
			}
			params.Metadata = map[string]string{
				"user_id": user.ID.String(),
			}
			//create stripe customer
			stripeCustomer, err := s.stripeClient.Customers.New(params)
			if err != nil {
				return nil, err
			}
			stripeLogin = &biz.UserLogin{
				UserId:        user.ID,
				LoginProvider: idp.StripeLoginProvider,
				ProviderKey:   stripeCustomer.ID,
			}
			err = s.um.AddLogin(ctx, user, []biz.UserLogin{*stripeLogin})
			if err != nil {
				return nil, err
			}
		}
		ret.StripeCustomerId = stripeLogin.ProviderKey
		ret.User = MapBizUserToApi(ctx, user, s.blob)
	} else if req.StripeCustomerId != nil {
		user, err := s.um.FindByLogin(ctx, idp.StripeLoginProvider, *req.StripeCustomerId)
		if err != nil {
			return nil, err
		}
		if user == nil {
			return nil, pb.ErrorUserNotFoundLocalized(ctx, nil, nil)
		}
		ret.StripeCustomerId = *req.StripeCustomerId
		ret.User = MapBizUserToApi(ctx, user, s.blob)
	}
	return ret, nil
}

// CheckUserTenant internal api for check user tenant
func (s *UserInternalService) CheckUserTenant(ctx context.Context, req *pb.CheckUserTenantRequest) (*pb.CheckUserTenantReply, error) {
	//check permission
	if err := kapi.ErrIfUntrusted(ctx, s.trust); err != nil {
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
