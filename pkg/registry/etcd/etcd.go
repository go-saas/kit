package etcd

import (
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	kregistry "github.com/go-saas/kit/pkg/registry"
	"github.com/goava/di"
	clientv3 "go.etcd.io/etcd/client/v3"
	"strings"
)

func init() {
	kregistry.Register("etcd", func(c *kregistry.Config, container *di.Container) (registry.Registrar, registry.Discovery, error) {
		var client *clientv3.Client
		if has, err := container.Has(&client); err == nil {
			// reuse client
			if !has {
				err := container.Provide(func() (*clientv3.Client, error) {
					ends := strings.Split(c.Endpoint, ",")
					return clientv3.New(clientv3.Config{
						Endpoints: ends,
					})
				})
				if err != nil {
					return nil, nil, err
				}
			} else {
				klog.Info("reuse etcd client")
			}
			err := container.Resolve(&client)
			if err != nil {
				return nil, nil, err
			}

		} else {
			return nil, nil, err
		}
		r := New(client)
		return r, r, nil
	})
}
