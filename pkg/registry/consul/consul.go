package consul

import (
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/registry"
	kregistry "github.com/go-saas/kit/pkg/registry"
	"github.com/hashicorp/consul/api"
)

func init() {
	kregistry.Register("consul", func(c *kregistry.Config) (registry.Registrar, registry.Discovery, error) {
		cfg := api.DefaultConfig()
		cfg.Address = c.Endpoint
		cli, err := api.NewClient(cfg)
		if err != nil {
			return nil, nil, err
		}
		r := consul.New(cli)
		return r, r, nil
	})
}
