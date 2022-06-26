package biz

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	mapstructure2 "github.com/go-saas/kit/pkg/mapstructure"
	"github.com/go-saas/saas/gorm"
	"github.com/go-saas/saas/seed"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"reflect"
	"strings"
)

type MenuSeed struct {
	dbProvider gorm.DbProvider
	menuRepo   MenuRepo
}

var _ seed.Contrib = (*MenuSeed)(nil)

//go:embed seed/menu.yaml
var menuData []byte

func NewMenuSeed(dbProvider gorm.DbProvider, menuRepo MenuRepo) *MenuSeed {
	return &MenuSeed{dbProvider: dbProvider, menuRepo: menuRepo}
}

func (m *MenuSeed) Seed(ctx context.Context, sCtx *seed.Context) error {
	if err := m.seedBytes(ctx, menuData); err != nil {
		return err
	}
	if seedPath, ok := sCtx.Extra[SeedPathKey]; ok {
		if path, ok := seedPath.(string); ok {
			b, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			if err := m.seedBytes(ctx, b); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *MenuSeed) seedBytes(ctx context.Context, data []byte) error {
	var v = make(map[string]interface{})
	dec := yaml.NewDecoder(bytes.NewReader(data))
	for {
		err := dec.Decode(v)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			} else {
				return err
			}
		}
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

				dbEntity.MergeWithPreservedFields(actual)
				if err := m.menuRepo.Update(ctx, pId, dbEntity, nil); err != nil {
					return err
				}
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
