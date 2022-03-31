package biz

import (
	"context"
	"github.com/goxiaoy/go-saas-kit/pkg/data"
	gorm2 "github.com/goxiaoy/go-saas-kit/pkg/gorm"
	"github.com/goxiaoy/go-saas-kit/pkg/query"
	v1 "github.com/goxiaoy/go-saas-kit/saas/api/tenant/v1"
	gg "gorm.io/gorm"
	"time"
)

type Tenant struct {
	gorm2.UIDBase
	//unique name. usually for domain name
	Name string `gorm:"column:name;index;size:255;"`
	//localed display name
	DisplayName string `gorm:"column:display_name;index;size:255;"`
	//region of this tenant
	Region    string `gorm:"column:region;index;size:255;"`
	Logo      string
	CreatedAt time.Time    `gorm:"column:created_at;index;"`
	UpdatedAt time.Time    `gorm:"column:updated_at;index;"`
	DeletedAt gg.DeletedAt `gorm:"column:deleted_at;index;"`

	//connection
	Conn []TenantConn `gorm:"foreignKey:TenantId"`
	//edition
	Features []TenantFeature `gorm:"foreignKey:TenantId"`
}

// TenantConn connection string info
type TenantConn struct {
	TenantId string `gorm:"column:tenant_id;primary_key;size:36;"`
	//key of connection string
	Key string `gorm:"column:key;primary_key;size:100;"`
	//connection string
	Value     string    `gorm:"column:value;size:1000;"`
	CreatedAt time.Time `gorm:"column:created_at;index;"`
	UpdatedAt time.Time `gorm:"column:updated_at;index;"`
}

type TenantFeature struct {
	TenantId string `gorm:"column:tenant_id;primary_key;size:36;"`
	//key of connection string
	Key string `gorm:"column:key;primary_key;size:100;"`
	//connection string
	Value     string    `gorm:"column:value;size:1000;"`
	CreatedAt time.Time `gorm:"column:created_at;index;"`
	UpdatedAt time.Time `gorm:"column:updated_at;index;"`
}

type TenantRepo interface {
	data.Repo[Tenant, string, v1.ListTenantRequest]
	FindByIdOrName(ctx context.Context, idOrName string) (*Tenant, error)
}

type TenantUseCase struct {
	repo TenantRepo
}

func NewTenantUserCase(repo TenantRepo) *TenantUseCase {
	return &TenantUseCase{repo: repo}
}

func (t TenantUseCase) FindByIdOrName(ctx context.Context, idOrName string) (*Tenant, error) {
	return t.repo.FindByIdOrName(ctx, idOrName)
}

func (t TenantUseCase) List(ctx context.Context, query *v1.ListTenantRequest) ([]*Tenant, error) {
	return t.repo.List(ctx, query)
}

func (t TenantUseCase) First(ctx context.Context, query *v1.ListTenantRequest) (*Tenant, error) {
	return t.repo.First(ctx, query)
}

func (t TenantUseCase) Count(ctx context.Context, query *v1.ListTenantRequest) (total int64, filtered int64, err error) {
	return t.repo.Count(ctx, query)
}

func (t TenantUseCase) Get(ctx context.Context, id string) (*Tenant, error) {
	return t.repo.Get(ctx, id)
}

func (t TenantUseCase) Create(ctx context.Context, entity *Tenant) error {
	// check duplicate
	dbEntity, err := t.repo.FindByIdOrName(ctx, entity.Name)
	if err != nil {
		return err
	}
	if dbEntity != nil {
		// duplicate
		return v1.ErrorDuplicateTenantName("%v is used", entity.Name)
	}
	return t.repo.Create(ctx, entity)
}

func (t TenantUseCase) Update(ctx context.Context, entity *Tenant, p query.Select) error {
	// check duplicate
	dbEntity, err := t.repo.FindByIdOrName(ctx, entity.Name)
	if err != nil {
		return err
	}
	if dbEntity != nil && dbEntity.ID != entity.ID {
		// duplicate
		return v1.ErrorDuplicateTenantName("%v is used", entity.Name)
	}
	return t.repo.Update(ctx, entity.ID.String(), entity, p)
}

func (t TenantUseCase) Delete(ctx context.Context, id string) error {
	return t.repo.Delete(ctx, id)
}
