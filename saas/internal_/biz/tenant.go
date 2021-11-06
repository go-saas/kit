package biz

import (
	"context"
	v1 "github.com/goxiaoy/go-saas-kit/saas/api/tenant/v1"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"time"
)

type Tenant struct {
	// unique id
	ID string
	//unique name. usually for domain name
	Name string
	//localed display name
	DisplayName string
	//region of this tenant. Useful for data storage location
	Region    string
	CreatedAt time.Time
	UpdatedAt time.Time
	//should apply soft delete
	DeletedAt *time.Time
	//connection
	Conn []TenantConn
	//edition
	Features []TenantFeature
}

type TenantRepo interface {
	FindByIdOrName(ctx context.Context, idOrName string) (*Tenant, error)
	List(ctx context.Context, query *v1.ListTenantRequest) ([]*Tenant, error)
	First(ctx context.Context, search string, query *v1.TenantFilter) (*Tenant, error)
	Count(ctx context.Context, search string, query *v1.TenantFilter) (total int64, filtered int64, err error)
	Get(ctx context.Context, id string) (*Tenant, error)
	Create(ctx context.Context, entity *Tenant) error
	Update(ctx context.Context, entity *Tenant, p *fieldmaskpb.FieldMask) error
	Delete(ctx context.Context, id string) error
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

func (t TenantUseCase) First(ctx context.Context, search string, query *v1.TenantFilter) (*Tenant, error) {
	return t.repo.First(ctx, search, query)
}

func (t TenantUseCase) Count(ctx context.Context, search string, query *v1.TenantFilter) (total int64, filtered int64, err error) {
	return t.repo.Count(ctx, search, query)
}

func (t TenantUseCase) Get(ctx context.Context, id string) (*Tenant, error) {
	return t.repo.Get(ctx, id)
}

func (t TenantUseCase) Create(ctx context.Context, entity *Tenant) error {
	// check duplicate
	dbEntity, err := t.repo.FindByIdOrName(ctx, entity.ID)
	if err != nil {
		return err
	}
	if dbEntity != nil {
		// duplicate
		return v1.ErrorDuplicateTenantName("%v is used", entity.Name)
	}
	return t.repo.Create(ctx, entity)
}

func (t TenantUseCase) Update(ctx context.Context, entity *Tenant, p *fieldmaskpb.FieldMask) error {
	// check duplicate
	dbEntity, err := t.repo.FindByIdOrName(ctx, entity.ID)
	if err != nil {
		return err
	}
	if dbEntity != nil && dbEntity.ID != entity.ID {
		// duplicate
		return v1.ErrorDuplicateTenantName("%v is used", entity.Name)
	}
	return t.repo.Update(ctx, entity, p)
}

func (t TenantUseCase) Delete(ctx context.Context, id string) error {
	return t.repo.Delete(ctx, id)
}
