package biz

import (
	"context"
	"github.com/go-saas/kit/pkg/apisix"
	"github.com/go-saas/kit/sys/private/conf"
	"github.com/go-saas/saas/seed"
	"github.com/hibiken/asynq"
	"math"
)

type ApisixSeed struct {
	Cfg       *conf.SysConf
	Client    *apisix.AdminClient
	JobClient *asynq.Client
}

func NewApisixSeed(cfg *conf.SysConf, client *apisix.AdminClient, jobClient *asynq.Client) *ApisixSeed {
	return &ApisixSeed{Cfg: cfg, Client: client, JobClient: jobClient}
}

func (a *ApisixSeed) Seed(ctx context.Context, sCtx *seed.Context) error {
	if len(sCtx.TenantId) != 0 || a.Cfg == nil || a.Cfg.Apisix == nil {
		return nil
	}
	//Put into background job
	_, err := a.JobClient.EnqueueContext(ctx, NewApisixMigrationTask(), asynq.MaxRetry(math.MaxInt))
	return err
}

var _ seed.Contrib = (*ApisixSeed)(nil)

func (a *ApisixSeed) Do() error {
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
