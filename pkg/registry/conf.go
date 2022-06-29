package registry

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-saas/kit/pkg/conf"
	"github.com/samber/lo"
)

type Conf struct {
	conf *conf.Services
	ctx  context.Context
}

var _ registry.Discovery = (*Conf)(nil)

func NewConf(conf *conf.Services) *Conf {
	return &Conf{conf: conf}
}

func (c *Conf) GetService(ctx context.Context, serviceName string) ([]*registry.ServiceInstance, error) {
	return mapToSrv(c.conf, serviceName)
}

func (c *Conf) Watch(ctx context.Context, serviceName string) (registry.Watcher, error) {
	return &confWatcher{ctx: ctx, conf: c.conf, serviceName: serviceName, first: true}, nil
}

type confWatcher struct {
	event       chan struct{}
	ctx         context.Context
	conf        *conf.Services
	serviceName string
	first       bool
}

func (e *confWatcher) Next() (services []*registry.ServiceInstance, err error) {
	if e.first {
		e.first = false
		return mapToSrv(e.conf, e.serviceName)
	}
	select {
	case <-e.ctx.Done():
		err = e.ctx.Err()
	case <-e.event:
	}
	return
}

func (e *confWatcher) Stop() error {
	return nil
}

var _ registry.Watcher = (*confWatcher)(nil)

func mapToSrv(c *conf.Services, serviceName string) ([]*registry.ServiceInstance, error) {
	if srvs, ok := c.Services[serviceName]; !ok {
		return nil, fmt.Errorf(" %v service not found", serviceName)
	} else {
		return lo.Map(srvs.List, func(s *conf.Service, _ int) *registry.ServiceInstance {
			return &registry.ServiceInstance{
				ID:        s.Id,
				Name:      s.Name,
				Version:   "",
				Metadata:  nil,
				Endpoints: s.Endpoints,
			}
		}), nil
	}
}
