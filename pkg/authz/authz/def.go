package authz

import (
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/samber/lo"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v3"
	"sort"
)

var (
	groups            []*PermissionDefGroup
	DefNotFoundReason = "PERMISSION_DEF_NOT_FOUND"
)

func init() {
	FindOrAddGroup(&PermissionDefGroup{
		Name: "internal",
		Side: PermissionAllowSide_BOTH,
		Def: []*PermissionDef{
			{
				Name:      "internal",
				Namespace: AnyNamespace,
				Action:    AnyAction.GetIdentity(),
			},
		},
		Internal: true,
	})
}

func FindOrAddGroup(group *PermissionDefGroup) *PermissionDefGroup {
	if len(group.Name) == 0 {
		panic(fmt.Errorf("group name required"))
	}
	find, ok := lo.Find(groups, func(t *PermissionDefGroup) bool {
		return t.Name == group.Name
	})
	if ok {
		return find
	}
	groups = append(groups, group)
	return group
}

func MustFindDef(namespace string, action Action) *PermissionDef {
	def, err := FindDef(namespace, action, false)
	if err != nil {
		panic(err)
	}
	return def
}

func FindDef(namespace string, action Action, publicOnly bool) (*PermissionDef, error) {
	def, ok := lo.Find(lo.FlatMap(groups, func(t *PermissionDefGroup, _ int) []*PermissionDef {
		return t.Def
	}), func(t *PermissionDef) bool {
		return t.Namespace == namespace && t.Action == action.GetIdentity()
	})
	if !ok || (publicOnly && def.Internal) {
		return nil, errors.New(400, DefNotFoundReason, fmt.Sprintf("permissin action %s in %s not defined", action.GetIdentity(), namespace))
	}
	return def, nil
}

func (x *PermissionDefGroup) AddDef(def *PermissionDef) {
	x.Def = append(x.Def, def)
}

func (x *PermissionDefGroup) NormalizeAndValidate() error {
	if len(x.Name) == 0 {
		return fmt.Errorf("group name required")
	}
	for _, def := range x.Def {
		if len(def.Name) == 0 {
			return fmt.Errorf("def under group %s name required", x.Name)
		}
		if x.Side != PermissionAllowSide_BOTH && x.Side != def.Side {
			if def.Side == PermissionAllowSide_BOTH {
				def.Side = x.Side
			} else {
				return fmt.Errorf("group %s has permission side %v, but try to add permission %s with %v side", x.Name, x.Side, def.Name, def.Side)
			}
		}
		if x.Internal {
			def.Internal = true
		}
	}
	return nil
}

func (x *PermissionDefGroup) Walk(isHost bool, publicOnly bool, f func(def *PermissionDef)) {

	var sortedDef []*PermissionDef
	sortedDef = append(sortedDef, x.Def...)
	sort.SliceStable(sortedDef, func(i, j int) bool {
		return sortedDef[i].Priority < sortedDef[j].Priority
	})

	for _, def := range sortedDef {
		if publicOnly && def.Internal {
			continue
		}
		if (def.Side == PermissionAllowSide_HOST_ONLY && !isHost) || (def.Side == PermissionAllowSide_TENANT_ONLY && isHost) {
			continue
		}
		f(def)
	}
}

func WalkGroups(isHost bool, publicOnly bool, f func(group *PermissionDefGroup)) {
	var sortedGroup []*PermissionDefGroup
	sortedGroup = append(sortedGroup, groups...)
	sort.SliceStable(sortedGroup, func(i, j int) bool {
		return sortedGroup[i].Priority < sortedGroup[j].Priority
	})
	for _, g := range sortedGroup {
		if publicOnly && g.Internal {
			continue
		}
		if (g.Side == PermissionAllowSide_HOST_ONLY && !isHost) || (g.Side == PermissionAllowSide_TENANT_ONLY && isHost) {
			continue
		}
		f(g)
	}
}

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
	groupConf := &PermissionConf{}
	err = opt.Unmarshal(b, groupConf)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	LoadFromConf(groupConf)
}

func LoadFromConf(groupConf *PermissionConf) {
	for _, group := range groupConf.Groups {
		//clone to clean previous def
		gg := proto.Clone(group).(*PermissionDefGroup)
		gg.Def = []*PermissionDef{}
		gg = FindOrAddGroup(gg)
		//merge def
		for _, def := range group.Def {
			gg.AddDef(def)
		}
		if err := gg.NormalizeAndValidate(); err != nil {
			panic(err)
		}
	}
}
