package casbin

import (
	"context"
	_ "embed"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/goxiaoy/go-saas/common"
	sgorm "github.com/goxiaoy/go-saas/gorm"
)

type EnforcerProvider struct {
	dbProvider sgorm.DbProvider
	dbKey      string
	m          model.Model
	logger     *klog.Helper
}

//go:embed model.conf
var modelStr string

func NewEnforcerProvider(logger klog.Logger, dbProvider sgorm.DbProvider, dbKey string) (*EnforcerProvider, error) {
	res := &EnforcerProvider{
		dbProvider: dbProvider,
		dbKey:      dbKey,
		logger:     klog.NewHelper(klog.With(logger, "module", "EnforcerProvider")),
	}
	if m, err := model.NewModelFromString(modelStr); err != nil {
		return nil, err
	} else {
		res.m = m
	}
	return res, nil
}

func (p *EnforcerProvider) Get(ctx context.Context) (*casbin.SyncedEnforcer, error) {

	db := p.dbProvider.Get(ctx, p.dbKey)

	var adapter *gormadapter.Adapter
	var err error
	adapter, err = gormadapter.NewAdapterByDB(db.WithContext(gormadapter.NewDisableAutoMigration(db.Statement.Context)))
	if err != nil {
		return nil, err
	}
	//every time reload policy by tenant
	tenantInfo, _ := common.FromCurrentTenant(ctx)
	filter := gormadapter.Filter{
		V4: []string{tenantInfo.GetId(), "*"},
	}
	//load filter policy
	err = adapter.LoadFilteredPolicy(p.m, filter)
	if err != nil {
		return nil, err
	}
	e, err := casbin.NewSyncedEnforcer(p.m, adapter)
	e.EnableAutoSave(true)

	return e, err

}
