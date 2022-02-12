package remote

import (
	"context"
	v1 "github.com/goxiaoy/go-saas-kit/saas/api/tenant/v1"
	"github.com/goxiaoy/go-saas/common"
)

type GrpcTenantStore struct {
	client v1.TenantServiceClient
}

var _ common.TenantStore = (*GrpcTenantStore)(nil)

func NewRemoteGrpcTenantStore(client v1.TenantServiceClient) common.TenantStore {
	return &GrpcTenantStore{client: client}
}

func (r *GrpcTenantStore) GetByNameOrId(ctx context.Context, nameOrId string) (*common.TenantConfig, error) {
	tenant, err := r.client.GetTenant(ctx, &v1.GetTenantRequest{IdOrName: nameOrId})
	if err != nil {
		return nil, err
	}
	ret := common.NewTenantConfig(tenant.Id, tenant.Name, tenant.Region)
	for _, conn := range tenant.Conn {
		ret.Conn[conn.Key] = conn.Value
	}
	return ret, nil
}