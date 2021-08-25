package seed

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	mapstructure2 "github.com/goxiaoy/go-saas-kit/pkg/mapstructure"
	"github.com/goxiaoy/go-saas-kit/user/internal_/biz"
	"github.com/goxiaoy/go-saas/seed"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
	"reflect"
)

//go:embed fake.yaml
var fake []byte

type Fake struct {
	um *biz.UserManager
}

const FakeSeedKey = "fake_seed"

func NewFake(um *biz.UserManager) *Fake {
	return &Fake{
		um,
	}
}

func (f Fake) Seed(ctx context.Context, sCtx *seed.Context) error {
	if v,ok:= sCtx.Extra[FakeSeedKey];ok {
		if b,ok:=v.(bool);ok{
			if b{
				var v = make(map[string]interface{})
				dec := yaml.NewDecoder(bytes.NewReader(fake))
				for dec.Decode(v) == nil {
					// find users
					if users,ok:= v["users"];ok{
						v := reflect.ValueOf(users)
						switch v.Kind() {
						case reflect.Slice:
							//fmt.Printf("slice: len=%d, %v\n", v.Len(), v.Interface())
							for i := 0; i < v.Len(); i++ {
								actual := &biz.User{}
								cfg := &mapstructure.DecoderConfig{
									DecodeHook: mapstructure.ComposeDecodeHookFunc(
										mapstructure2.StringToUUIDHookFunc(),
									),
									Metadata: nil,
									Result:   actual,
									TagName:  "json",
								}
								decoder, _ := mapstructure.NewDecoder(cfg)
								if err := decoder.Decode(v.Index(i).Interface());err!=nil{
									return err
								}
								dbUser,err:= f.um.FindByID(ctx,actual.ID.String())
								if err!=nil{
									return err
								}
								if dbUser!=nil{
									continue
								}
								if  err:=f.um.Create(ctx,actual);err!=nil{
									return err
								}
							}
						default:
							return errors.New("users should be array")
						}
					}
				}

			}
		}
	}

	return nil
}

