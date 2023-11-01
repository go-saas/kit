package api

import (
	"context"
	"github.com/go-saas/kit/pkg/api"
	v1 "github.com/go-saas/kit/saas/api/tenant/v1"
	"github.com/go-saas/saas"
)

type TenantStore struct {
	srv v1.TenantInternalServiceServer
}

var _ saas.TenantStore = (*TenantStore)(nil)

func NewTenantStore(srv v1.TenantInternalServiceServer) saas.TenantStore {
	return &TenantStore{srv: srv}
}

func (r *TenantStore) GetByNameOrId(ctx context.Context, nameOrId string) (*saas.TenantConfig, error) {
	//replace withe trusted environment to skip trusted check if in same process
	ctx = api.NewTrustedContext(ctx)
	tenant, err := r.srv.GetTenant(ctx, &v1.GetTenantRequest{IdOrName: nameOrId})
	if err != nil {
		return nil, err
	}
	pk := ""
	if tenant.PlanKey != nil {
		pk = *tenant.PlanKey
	}
	ret := saas.NewTenantConfig(tenant.Id, tenant.Name, tenant.Region, pk)
	for _, conn := range tenant.Conn {
		ret.Conn[conn.Key] = conn.Value
	}
	return ret, nil
}
