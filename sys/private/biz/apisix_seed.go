package biz

import (
	"context"
	"github.com/go-saas/kit/event"
	"github.com/go-saas/kit/pkg/apisix"
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
	"math"
)

type ApisixSeed struct {
	Cfg         *conf.SysConf
	Client      *apisix.AdminClient
	JobClient   *asynq.Client
	userSrv     v1.UserServiceServer
	eventSender event.Producer
}

func NewApisixSeed(cfg *conf.SysConf, client *apisix.AdminClient, jobClient *asynq.Client, userSrv v1.UserServiceServer, eventSender event.Producer) *ApisixSeed {
	return &ApisixSeed{Cfg: cfg, Client: client, JobClient: jobClient, userSrv: userSrv}
}

func (a *ApisixSeed) Seed(ctx context.Context, sCtx *seed.Context) error {
	if len(sCtx.TenantId) != 0 || a.Cfg == nil || a.Cfg.Apisix == nil {
		return nil
	}
	//Put into background job
	_, err := a.JobClient.EnqueueContext(ctx, NewApisixMigrationTask(), asynq.MaxRetry(math.MaxInt), asynq.Group(api.ServiceName))
	return err
}

var _ seed.Contrib = (*ApisixSeed)(nil)

func (a *ApisixSeed) Do(ctx context.Context) error {
	err := a.do(ctx)
	if err == nil {
		return nil
	}
	//send notification to admin users
	adminReply, err1 := a.userSrv.ListUsers(
		ctx,
		&v1.ListUsersRequest{
			PageSize: -1,
			Filter:   &v1.UserFilter{Roles: &v12.RoleFilter{Name: &query.StringFilterOperation{Eq: &wrapperspb.StringValue{Value: "admin"}}}},
			Fields:   &fieldmaskpb.FieldMask{Paths: []string{"id"}},
		},
	)
	if err1 != nil {
		return errors.Wrap(err, err1.Error())
	}
	adminIds := lo.Map(adminReply.Items, func(t *v1.User, _ int) string {
		return t.Id
	})
	notification := &v13.NotificationEvent{
		Title:   "Fail to Migrate Apisix Gateway",
		Desc:    err.Error(),
		UserIds: adminIds,
		Level:   v13.NotificationLevel_ERROR,
	}
	ee, _ := event.NewMessageFromProto(notification)
	err1 = a.eventSender.Send(ctx, ee)
	if err1 != nil {
		return errors.Wrap(err, err1.Error())
	}
	return err
}
func (a *ApisixSeed) do(_ context.Context) error {
	if a.Cfg.Apisix.Upstreams != nil {
		upstreams := a.Cfg.Apisix.Upstreams
		for id, upstream := range upstreams {
			if err := a.Client.PutUpstreamStruct(id, upstream); err != nil {
				return err
			}
		}
	}
	if a.Cfg.Apisix.GlobalRules != nil {
		rules := a.Cfg.Apisix.GlobalRules
		for id, rule := range rules {
			if err := a.Client.PutGlobalRules(id, rule); err != nil {
				return err
			}
		}
	}
	if a.Cfg.Apisix.Routes != nil {
		routes := a.Cfg.Apisix.Routes
		for id, route := range routes {
			if err := a.Client.PutRoute(id, route); err != nil {
				return err
			}
		}
	}
	if a.Cfg.Apisix.StreamRoutes != nil {
		routes := a.Cfg.Apisix.StreamRoutes
		for id, route := range routes {
			if err := a.Client.PutStreamRoute(id, route); err != nil {
				return err
			}
		}
	}
	return nil
}
