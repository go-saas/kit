package biz

import (
	"context"
	"github.com/google/uuid"
	"github.com/goxiaoy/go-saas-kit/pkg/data"
	"github.com/goxiaoy/go-saas-kit/pkg/gorm"
	"github.com/goxiaoy/go-saas-kit/pkg/query"
	v1 "github.com/goxiaoy/go-saas-kit/sys/api/menu/v1"
)

type Menu struct {
	gorm.UIDBase
	gorm.AuditedModel
	Name        string
	Desc        string
	Component   string
	Requirement []MenuPermissionRequirement `gorm:"foreignKey:MenuID"`
	Parent      string
	Props       data.JSONMap
	FullPath    string
	Priority    int32
	IgnoreAuth  bool
	Icon        string
}

type MenuPermissionRequirement struct {
	gorm.UIDBase
	MenuID    uuid.UUID `gorm:"type:char(36)"`
	Namespace string
	Resource  string
	Action    string
}

type MenuRepo interface {
	List(ctx context.Context, query *v1.ListMenuRequest) ([]*Menu, error)
	First(ctx context.Context, search string, query *v1.MenuFilter) (*Menu, error)
	Count(ctx context.Context, search string, query *v1.MenuFilter) (total int64, filtered int64, err error)
	Get(ctx context.Context, id string) (*Menu, error)
	Create(ctx context.Context, entity *Menu) error
	Update(ctx context.Context, entity *Menu, p query.Select) error
	Delete(ctx context.Context, id string) error
}
