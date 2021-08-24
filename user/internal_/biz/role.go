package biz

import (
	"context"
	"github.com/a8m/rql"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/goxiaoy/go-saas-kit/pkg/gorm"
	gorm2 "github.com/goxiaoy/go-saas/gorm"
	concurrency "github.com/goxiaoy/gorm-concurrency"
)

type Role struct {
	gorm.UIDBase
	concurrency.Version `gorm:"type:char(36)"`
	gorm2.MultiTenancy
	Name                string `json:"name" gorm:"index" rql:"filter"`
	NormalizedName      string `json:"normalized_name" gorm:"index"`
}

// RoleRepo crud role
type RoleRepo interface {
	List(ctx context.Context, query interface{}) ([]*Role, error)
	First(ctx context.Context, query interface{}) (*Role, error)
	FindByName(ctx context.Context, name string) (*Role, error)
	Count(ctx context.Context, query interface{}) (total int64, filtered int64, err error)
	Get(ctx context.Context, id string) (*Role, error)
	Create(ctx context.Context, role *Role) error
	Update(ctx context.Context, id string, role *Role, p rql.Select) error
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

func (r *RoleManager) First(ctx context.Context, query interface{}) (*Role, error) {
	return r.repo.First(ctx, query)
}

func (r *RoleManager) FindByName(ctx context.Context, name string) (*Role, error) {
	nn := r.lookupNormalizer.Name(name)
	return r.repo.FindByName(ctx, nn)
}

func (r *RoleManager) List(ctx context.Context, query interface{}) ([]*Role, error) {
	return r.repo.List(ctx, query)
}

func (r *RoleManager) Count(ctx context.Context, query interface{}) (total int64, filtered int64, err error) {
	return r.repo.Count(ctx, query)
}

func (r *RoleManager) Get(ctx context.Context, id string) (*Role, error) {
	return r.repo.Get(ctx, id)
}

func (r *RoleManager) Create(ctx context.Context, role *Role) error {
	nn := r.lookupNormalizer.Name(role.Name)
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

func (r *RoleManager) Update(ctx context.Context, id string, role *Role, p rql.Select) error {
	nn := r.lookupNormalizer.Name(role.Name)
	role.NormalizedName = nn
	dbRole, err := r.repo.FindByName(ctx, nn)
	if err != nil {
		return err
	}
	if dbRole != nil && dbRole.ID != role.ID {
		// duplicate
		return errors.Forbidden("NAME_DUPLICATE", "role name duplicate")
	}
	return r.repo.Update(ctx, id, role, p)
}

func (r *RoleManager) Delete(ctx context.Context, id string) error {
	return r.repo.Delete(ctx, id)
}
