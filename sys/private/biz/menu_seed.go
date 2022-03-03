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

func (m *MenuSeed) upsertMenus(ctx context.Context, parentId string, menus interface{}) error {
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
			actual.Parent = parentId
			//ensure create
			actual.Name = strings.ToLower(actual.Name)
			actual.IsPreserved = true
			dbEntity, err := m.menuRepo.FindByName(ctx, actual.Name)
			pId := ""
			if err != nil {
				return err
			}
			if dbEntity == nil {
				if err := m.menuRepo.Create(ctx, actual); err != nil {
					return err
				}
				pId = actual.ID.String()
			} else {
				pId = dbEntity.ID.String()
			}

			if children, ok := raw["children"]; ok {
				if err := m.upsertMenus(ctx, pId, children); err != nil {
					return err
				}
			}
		}
	default:
		return errors.New("users should be array")
	}
	return nil
}
