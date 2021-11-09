package api

import (
	"context"
	v1 "github.com/goxiaoy/go-saas-kit/saas/api/tenant/v1"
	"github.com/goxiaoy/go-saas/common"
)

type RemoteGrpcTenantStore struct {
	client v1.TenantServiceClient
}

var _ common.TenantStore = (*RemoteGrpcTenantStore)(nil)

func NewRemoteGrpcTenantStore(client v1.TenantServiceClient) common.TenantStore {
	return &RemoteGrpcTenantStore{client: client}
}

func (r *RemoteGrpcTenantStore) GetByNameOrId(ctx context.Context, nameOrId string) (*common.TenantConfig, error) {
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

type RemoteHttpTenantStore struct {
	client v1.TenantServiceHTTPClient
}

func NewRemoteHttpTenantStore(client v1.TenantServiceHTTPClient) common.TenantStore {
	return &RemoteHttpTenantStore{client: client}
}

func (r *RemoteHttpTenantStore) GetByNameOrId(ctx context.Context, nameOrId string) (*common.TenantConfig, error) {
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
