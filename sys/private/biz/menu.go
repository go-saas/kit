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
	gorm.UIDBase `json:",squash"`
	gorm.AuditedModel
	Name        string                      `json:"name"`
	Desc        string                      `json:"desc"`
	Component   string                      `json:"component"`
	Requirement []MenuPermissionRequirement `gorm:"foreignKey:MenuID;constraint:OnDelete:CASCADE;" json:"requirement"`
	Parent      string                      `json:"parent"`
	Props       data.JSONMap                `json:"props"`
	FullPath    string                      `json:"full_path"`
	Priority    int32                       `json:"priority"`
	IgnoreAuth  bool                        `json:"ignore_auth"`
	Icon        string                      `json:"icon"`
	Iframe      string                      `json:"iframe"`
	MicroApp    string                      `json:"micro_app"`
	Meta        data.JSONMap                `json:"meta"`
	Title       string                      `json:"title"`
	Path        string                      `json:"path"`
	Redirect    string                      `json:"redirect"`
	IsPreserved bool                        `json:"preserved"`
	HostOnly    bool                        `json:"host_only"`
}

type MenuPermissionRequirement struct {
	gorm.UIDBase `json:",squash"`
	MenuID       uuid.UUID `gorm:"type:char(36)" json:"menu_id"`
	Namespace    string    `json:"namespace"`
	Resource     string    `json:"resource"`
	Action       string    `json:"action"`
	HostOnly     bool      `json:"host_only"`
}

type MenuRepo interface {
	List(ctx context.Context, query *v1.ListMenuRequest) ([]*Menu, error)
	First(ctx context.Context, search string, query *v1.MenuFilter) (*Menu, error)
	FindByName(ctx context.Context, name string) (*Menu, error)
	Count(ctx context.Context, search string, query *v1.MenuFilter) (total int64, filtered int64, err error)
	Get(ctx context.Context, id string) (*Menu, error)
	Create(ctx context.Context, entity *Menu) error
	Update(ctx context.Context, entity *Menu, p query.Select) error
	Delete(ctx context.Context, id string) error
}
