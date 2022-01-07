package data

import (
	"context"
	"github.com/goxiaoy/go-saas-kit/saas/private/biz"
	"github.com/goxiaoy/go-saas/common"
)

// TenantStore query local to resolve tenant info
type TenantStore struct {
	tr biz.TenantRepo
}

func NewTenantStore(tr biz.TenantRepo) common.TenantStore {
	return &TenantStore{
		tr: tr,
	}
}

func (g TenantStore) GetByNameOrId(_ context.Context, nameOrId string) (*common.TenantConfig, error) {
	//change to host side
	newCtx := common.NewCurrentTenant(context.Background(), "", "")
	t, err := g.tr.FindByIdOrName(newCtx, nameOrId)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, common.ErrTenantNotFound
	}
	ret := common.NewTenantConfig(t.ID, t.Name, t.Region)
	for _, conn := range t.Conn {
		ret.Conn[conn.Key] = conn.Value
	}
	return ret, nil

}
