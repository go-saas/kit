package api

import (
	"context"
	"github.com/goxiaoy/go-saas-kit/pkg/api"
	v1 "github.com/goxiaoy/go-saas-kit/saas/api/tenant/v1"
	"github.com/goxiaoy/go-saas/common"
)

type TenantStore struct {
	srv v1.TenantInternalServiceServer
}

var _ common.TenantStore = (*TenantStore)(nil)

func NewTenantStore(srv v1.TenantInternalServiceServer) common.TenantStore {
	return &TenantStore{srv: srv}
}

func (r *TenantStore) GetByNameOrId(ctx context.Context, nameOrId string) (*common.TenantConfig, error) {
	//replace withe trusted environment to skip trusted check if in same process
	ctx = api.NewTrustedContext(ctx)
	tenant, err := r.srv.GetTenant(ctx, &v1.GetTenantRequest{IdOrName: nameOrId})
	if err != nil {
		return nil, err
	}
	ret := common.NewTenantConfig(tenant.Id, tenant.Name, tenant.Region)
	for _, conn := range tenant.Conn {
		ret.Conn[conn.Key] = conn.Value
	}
	return ret, nil
}
