package conf

import (
	"fmt"
	"github.com/go-saas/kit/event"
	"google.golang.org/protobuf/proto"
)

const (
	defaultKey = "default"
)

func (x *Services) GetServiceMergedDefault(name string) *Service {
	var res *Service
	if name != "" {
		res, _ = x.Services[name]
	}
	if def, ok := x.Services[defaultKey]; ok {
		if res == nil {
			res = def
		} else {
			c := proto.Clone(def).(*Service)
			proto.Merge(c, res)
			res = c
		}
	}

	return res
}

func (x *Endpoints) GetEventMergedDefault(name string) *event.Config {
	var res *event.Config
	if name != "" {
		res, _ = x.Events[name]
	}
	if def, ok := x.Events[defaultKey]; ok {
		if res == nil {
			res = def
		} else {
			c := proto.Clone(def).(*event.Config)
			proto.Merge(c, res)
			res = c
		}
	}
	if res == nil {
		panic(fmt.Sprintf("cannot resolve event %s or default", name))
	}
	return res
}

func (x *Endpoints) GetDatabaseMergedDefault(name string) *Database {
	var res *Database

	if name != "" {
		res, _ = x.Databases[name]
	}

	if def, ok := x.Databases[defaultKey]; ok {
		if res == nil {
			res = def
		} else {
			c := proto.Clone(def).(*Database)
			proto.Merge(c, res)
			res = c
		}
	}
	if res == nil {
		panic(fmt.Sprintf("cannot resolve database %s or default", name))
	}
	return res
}
