package casbin

import (
	"context"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/goxiaoy/go-saas/gorm"
)

type EnforcerProvider struct {
	dbProvider gorm.DbProvider
	modelPath  string
	modelStr   string
}

func NewEnforcerProvider(dbProvider gorm.DbProvider, modelPath string, modelStr string) *EnforcerProvider {
	return &EnforcerProvider{dbProvider: dbProvider}
}

func (p *EnforcerProvider) Get(ctx context.Context, key string) (*casbin.SyncedEnforcer, error) {
	db := p.dbProvider.Get(ctx, key)
	a, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}

	if p.modelStr != "" {
		if m, err := model.NewModelFromString(p.modelStr); err != nil {
			return nil, err
		} else {
			e, err :=casbin.NewSyncedEnforcer(m, a)
			return e,err
		}

	}
	if p.modelPath != "" {
		e, err := casbin.NewSyncedEnforcer(p.modelPath, a)
		return e, err
	}
	panic("modelPath or modelStr can not be empty")
}
