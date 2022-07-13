package biz

import (
	"context"
	"github.com/go-saas/kit/pkg/apisix"
	"github.com/go-saas/kit/sys/private/conf"
	"github.com/go-saas/saas/seed"
)

type ApisixSeed struct {
	Cfg    *conf.SysConf
	Client *apisix.AdminClient
}

var _ seed.Contrib = (*ApisixSeed)(nil)

func (a *ApisixSeed) Seed(ctx context.Context, sCtx *seed.Context) error {
	if len(sCtx.TenantId) != 0 || a.Cfg == nil || a.Cfg.Apisix == nil {
		return nil
	}
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
	return nil
}
