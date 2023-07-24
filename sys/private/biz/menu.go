package biz

import (
	"context"
	"github.com/go-saas/kit/pkg/data"
	"github.com/go-saas/kit/pkg/gorm"
	v1 "github.com/go-saas/kit/sys/api/menu/v1"
	"github.com/google/uuid"
)

type Menu struct {
	gorm.UIDBase `json:",squash"`
	gorm.AuditedModel
	Name              string                      `json:"name"`
	Desc              string                      `json:"desc"`
	Component         string                      `json:"component"`
	Requirement       []MenuPermissionRequirement `gorm:"foreignKey:MenuID;constraint:OnDelete:CASCADE;" json:"requirement"`
	Parent            string                      `json:"parent"`
	Props             data.JSONMap                `json:"props"`
	FullPath          string                      `json:"full_path"`
	Priority          int32                       `json:"priority"`
	IgnoreAuth        bool                        `json:"ignore_auth"`
	Icon              string                      `json:"icon"`
	Iframe            string                      `json:"iframe"`
	MicroApp          string                      `json:"micro_app"`
	MicroAppDev       string                      `json:"micro_app_dev"`
	MicroAppName      string                      `json:"micro_app_name"`
	MicroAppBaseRoute string                      `json:"micro_app_base_route"`
	Meta              data.JSONMap                `json:"meta"`
	Title             string                      `json:"title"`
	Path              string                      `json:"path"`
	Redirect          string                      `json:"redirect"`
	HideInMenu        bool                        `json:"hide_in_menu"`
	IsPreserved       bool                        `json:"is_preserved"`
}

type MenuPermissionRequirement struct {
	gorm.UIDBase `json:",squash"`
	MenuID       uuid.UUID `gorm:"type:char(36)" json:"menu_id"`
	Namespace    string    `json:"namespace"`
	Resource     string    `json:"resource"`
	Action       string    `json:"action"`
}

type MenuRepo interface {
	data.Repo[Menu, string, *v1.ListMenuRequest]
	FindByName(ctx context.Context, name string) (*Menu, error)
}

func (m *Menu) MergeWithPreservedFields(p *Menu) {
	m.Name = p.Name
	m.Component = p.Component
	m.Requirement = p.Requirement
	m.Iframe = p.Iframe
	m.Parent = p.Parent
	m.FullPath = p.FullPath
	m.Redirect = p.Redirect
	m.IsPreserved = p.IsPreserved
	m.Path = p.Path
	m.Priority = p.Priority
	m.Title = p.Title
	m.Icon = p.Icon
	m.HideInMenu = p.HideInMenu
	m.MicroApp = p.MicroApp
	m.MicroAppName = p.MicroAppName
	m.MicroAppBaseRoute = p.MicroAppBaseRoute
}
