package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/goxiaoy/go-saas-kit/pkg/gorm"
	"github.com/goxiaoy/go-saas-kit/pkg/query"
	v12 "github.com/goxiaoy/go-saas-kit/user/api/role/v1"
	gorm2 "github.com/goxiaoy/go-saas/gorm"
	concurrency "github.com/goxiaoy/gorm-concurrency"
)

type Role struct {
	gorm.UIDBase
	concurrency.Version `gorm:"type:char(36)"`
	gorm.AuditedModel
	gorm2.MultiTenancy
	Name           string `json:"name" gorm:"index"`
	NormalizedName string `json:"normalized_name" gorm:"index"`
	IsPreserved    bool   `json:"is_preserved"`
}

// RoleRepo crud role
type RoleRepo interface {
	List(ctx context.Context, query *v12.ListRolesRequest) ([]*Role, error)
	First(ctx context.Context, query *v12.RoleFilter) (*Role, error)
	FindById(ctx context.Context, id string) (*Role, error)
	FindByName(ctx context.Context, name string) (*Role, error)
	Count(ctx context.Context, query *v12.RoleFilter) (total int64, filtered int64, err error)
	Get(ctx context.Context, id string) (*Role, error)
	Create(ctx context.Context, role *Role) error
	Update(ctx context.Context, id string, role *Role, p query.Select) error
	Delete(ctx context.Context, id string) error
}

var _ RoleRepo = (*RoleManager)(nil)

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

func (r *RoleManager) First(ctx context.Context, query *v12.RoleFilter) (*Role, error) {
	return r.repo.First(ctx, query)
}

func (r *RoleManager) FindByName(ctx context.Context, name string) (*Role, error) {
	nn, err := r.lookupNormalizer.Name(name)
	if err != nil {
		return nil, err
	}
	return r.repo.FindByName(ctx, nn)
}

func (r *RoleManager) FindById(ctx context.Context, id string) (*Role, error) {
	return r.repo.FindById(ctx, id)
}
func (r *RoleManager) List(ctx context.Context, query *v12.ListRolesRequest) ([]*Role, error) {
	return r.repo.List(ctx, query)
}

func (r *RoleManager) Count(ctx context.Context, query *v12.RoleFilter) (total int64, filtered int64, err error) {
	return r.repo.Count(ctx, query)
}

func (r *RoleManager) Get(ctx context.Context, id string) (*Role, error) {
	return r.repo.Get(ctx, id)
}

func (r *RoleManager) Create(ctx context.Context, role *Role) error {
	nn, err := r.lookupNormalizer.Name(role.Name)
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
		return errors.Forbidden("NAME_DUPLICATE", "role name duplicate")
	}
	role.NormalizedName = nn
	return r.repo.Create(ctx, role)
}

func (r *RoleManager) Update(ctx context.Context, id string, role *Role, p query.Select) error {
	nn, err := r.lookupNormalizer.Name(role.Name)
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
		return errors.Forbidden("NAME_DUPLICATE", "role name duplicate")
	}
	if role.IsPreserved {
		return v12.ErrorRolePreserved("")
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
		return v12.ErrorRolePreserved("")
	}
	return r.repo.Delete(ctx, id)
}
