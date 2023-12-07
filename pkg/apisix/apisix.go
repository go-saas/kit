package apisix

import (
	"encoding/json"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/encoding/protojson"
	"gopkg.in/yaml.v3"
)

var (
	modules []*Module
)

func LoadFromYaml(data []byte) {
	m := make(map[string]interface{})

	err := yaml.Unmarshal(data, &m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	b, err := json.Marshal(m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	opt := protojson.UnmarshalOptions{}
	module := &Module{}
	err = opt.Unmarshal(b, module)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	modules = append(modules, module)
}

func WalkModules(f func(module *Module) error) error {
	for _, module := range modules {
		if err := f(module); err != nil {
			return err
		}
	}
	return nil
}
