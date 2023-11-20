package casbin

import (
	"context"
	_ "embed"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-saas/saas"

	sgorm "github.com/go-saas/saas/gorm"
)

type EnforcerProvider struct {
	dbProvider sgorm.DbProvider
	dbKey      string
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
	if _, err := model.NewModelFromString(modelStr); err != nil {
		return nil, err
	} else {

	}
	return res, nil
}

func (p *EnforcerProvider) Get(ctx context.Context) (*casbin.SyncedEnforcer, error) {

	db := p.dbProvider.Get(ctx, p.dbKey)

	var adapter *gormadapter.Adapter
	var err error

	gormadapter.TurnOffAutoMigrate(db)
	adapter, err = gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}
	//every time reload policy by tenant
	ti, _ := saas.FromCurrentTenant(ctx)
	filter := gormadapter.Filter{
		V4: []string{ti.GetId(), "*"},
	}
	//model is not concurrency safe
	m, _ := model.NewModelFromString(modelStr)
	//load filter policy
	err = adapter.LoadFilteredPolicy(m, filter)
	if err != nil {
		return nil, err
	}
	e, err := casbin.NewSyncedEnforcer(m, adapter)
	e.EnableAutoSave(true)

	return e, err

}
