package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-saas/kit/pkg/data"
	kitgorm "github.com/go-saas/kit/pkg/gorm"
	"github.com/go-saas/kit/pkg/query"
	v12 "github.com/go-saas/kit/user/api/role/v1"
	sgorm "github.com/go-saas/saas/gorm"
	concurrency "github.com/goxiaoy/gorm-concurrency/v2"
)

type Role struct {
	kitgorm.UIDBase
	concurrency.HasVersion `gorm:"type:char(36)"`
	kitgorm.AuditedModel
	TenantId       sgorm.HasTenant `gorm:"index:idx_tenant_role,unique"`
	Name           string          `json:"name" gorm:"size:200"`
	NormalizedName string          `gorm:"size:200;index:idx_tenant_role,unique" json:"normalized_name" `
	IsPreserved    bool            `json:"is_preserved"`
}

// RoleRepo crud role
type RoleRepo interface {
	data.Repo[Role, string, *v12.ListRolesRequest]
	FindByName(ctx context.Context, name string) (*Role, error)
}

//var _ RoleRepo = (*RoleManager)(nil)

type RoleManager struct {
	repo             RoleRepo
	lookupNormalizer LookupNormalizer
}

func NewRoleManager(repo RoleRepo, lookupNormalizer LookupNormalizer) *RoleManager {
	return &RoleManager{
		repo:             repo,
		lookupNormalizer: lookupNormalizer,
	}
}

func (r *RoleManager) First(ctx context.Context, query *v12.ListRolesRequest) (*Role, error) {
	return r.repo.First(ctx, query)
}

func (r *RoleManager) FindByName(ctx context.Context, name string) (*Role, error) {
	nn, err := r.lookupNormalizer.Name(ctx, name)
	if err != nil {
		return nil, err
	}
	return r.repo.FindByName(ctx, nn)
}

func (r *RoleManager) List(ctx context.Context, query *v12.ListRolesRequest) ([]*Role, error) {
	return r.repo.List(ctx, query)
}

func (r *RoleManager) Count(ctx context.Context, query *v12.ListRolesRequest) (total int64, filtered int64, err error) {
	return r.repo.Count(ctx, query)
}

func (r *RoleManager) Get(ctx context.Context, id string) (*Role, error) {
	return r.repo.Get(ctx, id)
}

func (r *RoleManager) Create(ctx context.Context, role *Role) error {
	nn, err := r.lookupNormalizer.Name(ctx, role.Name)
	if err != nil {
		return err
	}
	// check duplicate
	dbRole, err := r.repo.FindByName(ctx, nn)
	if err != nil {
		return err
	}
	if dbRole != nil {
		// duplicate
		return v12.ErrorRoleNameDuplicateLocalized(ctx, nil, nil)
	}
	role.NormalizedName = nn
	return r.repo.Create(ctx, role)
}

func (r *RoleManager) Update(ctx context.Context, id string, role *Role, p query.Select) error {
	nn, err := r.lookupNormalizer.Name(ctx, role.Name)
	if err != nil {
		return err
	}
	role.NormalizedName = nn
	dbRole, err := r.repo.FindByName(ctx, nn)
	if err != nil {
		return err
	}
	if dbRole != nil && dbRole.ID != role.ID {
		// duplicate
		return v12.ErrorRoleNameDuplicateLocalized(ctx, nil, nil)
	}
	if role.IsPreserved {
		return v12.ErrorRolePreservedLocalized(ctx, nil, nil)
	}
	return r.repo.Update(ctx, id, role, p)
}

func (r *RoleManager) Delete(ctx context.Context, id string) error {
	role, err := r.Get(ctx, id)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.NotFound("", "")
	}
	if role.IsPreserved {
		return v12.ErrorRolePreservedLocalized(ctx, nil, nil)
	}
	return r.repo.Delete(ctx, id)
}
