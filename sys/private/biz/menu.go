package biz

import (
	"context"
	"github.com/google/uuid"
	"github.com/goxiaoy/go-saas-kit/pkg/data"
	"github.com/goxiaoy/go-saas-kit/pkg/gorm"
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
	IsPreserved bool                        `json:"is_preserved"`
}

type MenuPermissionRequirement struct {
	gorm.UIDBase `json:",squash"`
	MenuID       uuid.UUID `gorm:"type:char(36)" json:"menu_id"`
	Namespace    string    `json:"namespace"`
	Resource     string    `json:"resource"`
	Action       string    `json:"action"`
}

type MenuRepo interface {
	data.Repo[Menu, string, v1.ListMenuRequest]
	FindByName(ctx context.Context, name string) (*Menu, error)
}

func (m *Menu) MergeWithPreservedFields(p *Menu) {
	m.Name = p.Name
	m.Component = p.Component
	m.Requirement = p.Requirement
	m.MicroApp = p.MicroApp
	m.Iframe = p.Iframe
	m.Parent = p.Parent
	m.FullPath = p.FullPath
	m.Redirect = p.Redirect
	m.IsPreserved = p.IsPreserved
}
