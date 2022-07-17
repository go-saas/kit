package consul

import (
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	kregistry "github.com/go-saas/kit/pkg/registry"
	"github.com/goava/di"
	"github.com/hashicorp/consul/api"
)

func init() {
	kregistry.Register("consul", func(c *kregistry.Config, container *di.Container) (registry.Registrar, registry.Discovery, error) {
		var client *api.Client
		if has, err := container.Has(&client); err == nil {
			if !has {
				err := container.Provide(func() (*api.Client, error) {
					cfg := api.DefaultConfig()
					cfg.Address = c.Endpoint
					return api.NewClient(cfg)
				})
				if err != nil {
					return nil, nil, err
				}
			}
			err := container.Resolve(&client)
			if err != nil {
				return nil, nil, err
			} else {
				klog.Info("reuse consul client")
			}
		} else {
			return nil, nil, err
		}

		r := consul.New(client)
		return r, r, nil
	})
}
