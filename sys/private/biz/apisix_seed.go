package biz

import (
	"context"
	"github.com/go-saas/kit/event"
	"github.com/go-saas/kit/pkg/apisix"
	"github.com/go-saas/kit/pkg/authz/authz"
	"github.com/go-saas/kit/pkg/query"
	v13 "github.com/go-saas/kit/realtime/event/v1"
	"github.com/go-saas/kit/sys/api"
	"github.com/go-saas/kit/sys/private/conf"
	v12 "github.com/go-saas/kit/user/api/role/v1"
	v1 "github.com/go-saas/kit/user/api/user/v1"
	"github.com/go-saas/saas/seed"
	"github.com/hibiken/asynq"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type ApisixSeed struct {
	Cfg         *conf.SysConf
	Client      *apisix.AdminClient
	JobClient   *asynq.Client
	userSrv     v1.UserServiceServer
	eventSender event.Producer
}

func NewApisixSeed(cfg *conf.SysConf, client *apisix.AdminClient, jobClient *asynq.Client, userSrv v1.UserServiceServer, eventSender event.Producer) *ApisixSeed {
	return &ApisixSeed{Cfg: cfg, Client: client, JobClient: jobClient, userSrv: userSrv, eventSender: eventSender}
}

func (a *ApisixSeed) Seed(ctx context.Context, sCtx *seed.Context) error {
	if len(sCtx.TenantId) != 0 || a.Cfg == nil || a.Cfg.Apisix == nil {
		return nil
	}
	//Put into background job
	_, err := a.JobClient.EnqueueContext(ctx, NewApisixMigrationTask(), asynq.MaxRetry(1000), asynq.Group(api.ServiceName))
	return err
}

var _ seed.Contrib = (*ApisixSeed)(nil)

func (a *ApisixSeed) Do(ctx context.Context) error {
	err := a.do(ctx)
	var notification *v13.NotificationEvent
	formatErr := func(err2 error) error {
		if err == nil {
			return err2
		} else {
			return errors.Wrap(err2, err.Error())
		}
	}
	if err == nil {
		notification = &v13.NotificationEvent{
			Title: "Migrate Apisix Gateway Successfully",
			Level: v13.NotificationLevel_INFO,
		}
	} else {
		notification = &v13.NotificationEvent{
			Title: "Fail to Migrate Apisix Gateway",
			Desc:  err.Error(),
			Level: v13.NotificationLevel_ERROR,
		}
	}
	//send notification to admin users
	tempCtx := authz.NewAlwaysAuthorizationContext(ctx, true)
	adminReply, err1 := a.userSrv.ListUsers(
		tempCtx,
		&v1.ListUsersRequest{
			PageSize: -1,
			Filter:   &v1.UserFilter{Roles: &v12.RoleFilter{Name: &query.StringFilterOperation{Eq: &wrapperspb.StringValue{Value: "admin"}}}},
			Fields:   &fieldmaskpb.FieldMask{Paths: []string{"id"}},
		},
	)
	if err1 != nil {
		return formatErr(err1)
	}
	adminIds := lo.Map(adminReply.Items, func(t *v1.User, _ int) string {
		return t.Id
	})
	notification.UserIds = adminIds
	ee, _ := event.NewMessageFromProto(notification)
	err1 = a.eventSender.Send(ctx, ee)
	if err1 != nil {
		return formatErr(err1)
	}
	return err
}

func (a *ApisixSeed) do(_ context.Context) error {
	if err := apisix.WalkModules(a.doModule); err != nil {
		return err
	}
	if a.Cfg.Apisix != nil {
		for _, module := range a.Cfg.Apisix.Modules {
			if err := a.doModule(module); err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *ApisixSeed) doModule(module *apisix.Module) error {
	if module.Upstreams != nil {
		for id, upstream := range module.Upstreams {
			if err := a.Client.PutUpstreamStruct(id, upstream); err != nil {
				return err
			}
		}
	}
	if module.GlobalRules != nil {
		for id, rule := range module.GlobalRules {
			if err := a.Client.PutGlobalRules(id, rule); err != nil {
				return err
			}
		}
	}
	if module.Routes != nil {
		for id, route := range module.Routes {
			if err := a.Client.PutRoute(id, route); err != nil {
				return err
			}
		}
	}
	if module.StreamRoutes != nil {
		for id, route := range module.StreamRoutes {
			if err := a.Client.PutStreamRoute(id, route); err != nil {
				return err
			}
		}
	}
	return nil
}
