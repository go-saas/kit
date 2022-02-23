package biz

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	mapstructure2 "github.com/goxiaoy/go-saas-kit/pkg/mapstructure"
	"github.com/goxiaoy/go-saas/gorm"
	"github.com/goxiaoy/go-saas/seed"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
	"reflect"
	"strings"
)

type MenuSeed struct {
	dbProvider gorm.DbProvider
	menuRepo   MenuRepo
}

var _ seed.Contributor = (*MenuSeed)(nil)

//go:embed seed/menu.yaml
var menuData []byte

func NewMenuSeed(dbProvider gorm.DbProvider, menuRepo MenuRepo) *MenuSeed {
	return &MenuSeed{dbProvider: dbProvider, menuRepo: menuRepo}
}

func (m *MenuSeed) Seed(ctx context.Context, sCtx *seed.Context) error {
	var v = make(map[string]interface{})
	dec := yaml.NewDecoder(bytes.NewReader(menuData))
	for dec.Decode(v) == nil {
		if menus, ok := v["menus"]; ok {
			err := m.upsertMenus(ctx, "", menus)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *MenuSeed) upsertMenus(ctx context.Context, parent string, menus interface{}) error {
	v := reflect.ValueOf(menus)
	switch v.Kind() {
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			raw := v.Index(i).Interface().(map[string]interface{})
			actual := &Menu{}
			cfg := &mapstructure.DecoderConfig{
				DecodeHook: mapstructure.ComposeDecodeHookFunc(
					mapstructure2.StringToUUIDHookFunc(),
					mapstructure2.JsonToJsonMapHookFunc(),
				),
				Metadata: nil,
				Result:   actual,
				TagName:  "json",
			}
			decoder, _ := mapstructure.NewDecoder(cfg)
			if err := decoder.Decode(raw); err != nil {
				return err
			}

			//find by name
			if actual.Name == "" {
				return errors.New("menu name is required")
			}
			actual.Parent = parent
			//ensure create
			actual.Name = strings.ToLower(actual.Name)
			dbEntity, err := m.menuRepo.FindByName(ctx, actual.Name)
			if err != nil {
				return err
			}
			if dbEntity == nil {
				if err := m.menuRepo.Create(ctx, actual); err != nil {
					return err
				}
			}

			if children, ok := raw["children"]; ok {
				if err := m.upsertMenus(ctx, actual.Name, children); err != nil {
					return err
				}
			}
		}
	default:
		return errors.New("users should be array")
	}
	return nil
}
