package remote

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	v1 "github.com/goxiaoy/go-saas-kit/saas/api/tenant/v1"
	"github.com/goxiaoy/go-saas/common"
)

type GrpcTenantStore struct {
	client v1.TenantServiceClient
}

var _ common.TenantStore = (*GrpcTenantStore)(nil)

func NewRemoteGrpcTenantStore(client v1.TenantServiceClient) common.TenantStore {
	return common.NewCachedTenantStore(&GrpcTenantStore{client: client})
}

func (r *GrpcTenantStore) GetByNameOrId(ctx context.Context, nameOrId string) (*common.TenantConfig, error) {
	tenant, err := r.client.GetTenantInternal(ctx, &v1.GetTenantRequest{IdOrName: nameOrId})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, common.ErrTenantNotFound
		}
		return nil, err
	}
	ret := common.NewTenantConfig(tenant.Id, tenant.Name, tenant.Region)
	for _, conn := range tenant.Conn {
		ret.Conn[conn.Key] = conn.Value
	}
	return ret, nil
}
