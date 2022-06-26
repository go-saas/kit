package biz

import (
	"context"
	"fmt"
	"github.com/go-saas/kit/saas/private/conf"
)

type ConnStrGenerator interface {
	//Generate connection string for tenant before creation
	Generate(ctx context.Context, tenant *Tenant) ([]TenantConn, error)
}

type ConfigConnStrGenerator struct {
	saasConf *conf.SaasConf
}

func NewConfigConnStrGenerator(saasConf *conf.SaasConf) ConnStrGenerator {
	return &ConfigConnStrGenerator{saasConf: saasConf}
}

func (c *ConfigConnStrGenerator) Generate(ctx context.Context, tenant *Tenant) ([]TenantConn, error) {
	var res []TenantConn
	if c.saasConf != nil {
		for _, template := range c.saasConf.Database {
			res = append(res, TenantConn{Key: template.Name, Value: fmt.Sprintf(template.Template, tenant.ID.String())})
		}
	}
	return res, nil
}
