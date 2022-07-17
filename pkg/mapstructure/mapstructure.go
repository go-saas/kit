package mapstructure

import (
	data2 "github.com/go-saas/kit/pkg/data"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"reflect"
)

func StringToUUIDHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(uuid.UUID{}) {
			return data, nil
		}

		return uuid.Parse(data.(string))
	}
}

func JsonToJsonMapHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {

		if v, ok := data.(map[string]interface{}); ok {
			return data2.JSONMap(v), nil
		}
		return data, nil
	}
}
