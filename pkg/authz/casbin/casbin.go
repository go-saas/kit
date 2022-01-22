package casbin

import (
	"context"
	_ "embed"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/goxiaoy/go-saas/common"
	"github.com/goxiaoy/go-saas/gorm"
)

type EnforcerProvider struct {
	dbProvider gorm.DbProvider
	dbKey      string
}

//go:embed model.conf
var modelStr string

func NewEnforcerProvider(dbProvider gorm.DbProvider, dbKey string) *EnforcerProvider {
	return &EnforcerProvider{dbProvider: dbProvider, dbKey: dbKey}
}

func (p *EnforcerProvider) Get(ctx context.Context) (*casbin.SyncedEnforcer, error) {
	db := p.dbProvider.Get(ctx, p.dbKey)
	a, err := gormadapter.NewAdapterByDB(db)

	if err != nil {
		return nil, err
	}
	tenantInfo := common.FromCurrentTenant(ctx)
	filter := gormadapter.Filter{
		//TODO host side?
		V4: []string{tenantInfo.GetId(), "*"},
	}
	if m, err := model.NewModelFromString(modelStr); err != nil {
		return nil, err
	} else {
		err := a.LoadFilteredPolicy(m, filter)
		if err != nil {
			return nil, err
		}
		e, err := casbin.NewSyncedEnforcer(m, a)
		e.EnableAutoSave(true)
		return e, err
	}

	panic("modelPath or modelStr can not be empty")
}
