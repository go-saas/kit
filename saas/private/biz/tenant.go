package biz

import (
	"context"
	"github.com/go-saas/kit/pkg/data"
	kitgorm "github.com/go-saas/kit/pkg/gorm"
	"github.com/go-saas/kit/pkg/query"
	v1 "github.com/go-saas/kit/saas/api/tenant/v1"
	"github.com/google/uuid"
	concurrency "github.com/goxiaoy/gorm-concurrency/v2"
	"gorm.io/gorm"
	"time"
)

type Tenant struct {
	kitgorm.UIDBase
	kitgorm.AggRoot
	concurrency.HasVersion
	//unique name. usually for domain name
	Name string `gorm:"column:name;index;size:255;"`
	//localed display name
	DisplayName string `gorm:"column:display_name;index;size:255;"`
	//region of this tenant
	Region     string `gorm:"column:region;index;size:255;"`
	Logo       string
	CreatedAt  time.Time      `gorm:"column:created_at;index;"`
	UpdatedAt  time.Time      `gorm:"column:updated_at;index;"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index;"`
	SeparateDb bool
	//connection
	Conn []TenantConn `gorm:"foreignKey:TenantId"`

	//Plan
	PlanKey              *string
	Plan                 *Plan `gorm:"foreignKey:PlanKey"`
	ActiveSubscriptionID *string

	Features []TenantFeature `gorm:"foreignKey:TenantId"`
	Extra    data.JSONMap
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
	Ready     bool      `gorm:"column:ready"`
}

type TenantFeature struct {
	kitgorm.AuditedModel
	TenantId string     `gorm:"column:tenant_id;primary_key;size:36;"`
	Key      string     `gorm:"column:key;primary_key;size:100;"`
	Value    data.Value `gorm:"embedded"`
}

type TenantRepo interface {
	data.Repo[Tenant, string, *v1.ListTenantRequest]
	FindByIdOrName(ctx context.Context, idOrName string) (*Tenant, error)
}

type TenantUseCase struct {
	TenantRepo
	connStrGenerator ConnStrGenerator
}

func NewTenantUserCase(repo TenantRepo, connStrGenerator ConnStrGenerator) *TenantUseCase {
	return &TenantUseCase{TenantRepo: repo, connStrGenerator: connStrGenerator}
}

func (t *TenantUseCase) Create(ctx context.Context, entity *Tenant) error {
	// check duplicate
	dbEntity, err := t.TenantRepo.FindByIdOrName(ctx, entity.Name)
	if err != nil {
		return err
	}
	if dbEntity != nil {
		// duplicate
		return v1.ErrorDuplicateTenantNameLocalized(ctx, map[string]interface{}{"name": entity.Name}, nil)
	}
	//ensure id generate
	if entity.UIDBase.ID == uuid.Nil {
		entity.UIDBase.ID = uuid.New()
	}

	if entity.SeparateDb {
		conn, err := t.connStrGenerator.Generate(ctx, entity)
		if err != nil {
			return err
		}
		entity.Conn = conn
	}

	if err := t.TenantRepo.Create(ctx, entity); err != nil {
		return err
	}
	return nil
}

func (t *TenantUseCase) Update(ctx context.Context, entity *Tenant, p query.Select) error {
	// check duplicate
	dbEntity, err := t.TenantRepo.FindByIdOrName(ctx, entity.Name)
	if err != nil {
		return err
	}
	if dbEntity != nil && dbEntity.ID != entity.ID {
		// duplicate
		return v1.ErrorDuplicateTenantNameLocalized(ctx, map[string]interface{}{"name": entity.Name}, nil)
	}
	return t.TenantRepo.Update(ctx, entity.ID.String(), entity, p)
}
