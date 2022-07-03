package etcd

import (
	etcd "github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/registry"
	kregistry "github.com/go-saas/kit/pkg/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
	"strings"
)

func init() {
	kregistry.Register("etcd", func(c *kregistry.Config) (registry.Registrar, registry.Discovery, error) {
		ends := strings.Split(c.Endpoint, ",")
		cli, err := clientv3.New(clientv3.Config{
			Endpoints: ends,
		})
		if err != nil {
			return nil, nil, err
		}
		r := etcd.New(cli)
		return r, r, nil
	})
}
