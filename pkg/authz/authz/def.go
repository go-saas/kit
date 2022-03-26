package authz

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/samber/lo"
	"sort"
	"sync"
)

var (
	//TODO priority queue
	groups []*PermissionDefGroup
	lock   sync.Mutex
)

type PermissionDefGroup struct {
	DisplayName string
	Side        PermissionSide
	Priority    int
	//TODO priority queue
	def   []PermissionDef
	Extra map[string]interface{}
}

func NewPermissionDefGroup(displayName string, side PermissionSide, priority int, extra ...map[string]interface{}) *PermissionDefGroup {
	res := &PermissionDefGroup{DisplayName: displayName, Side: side, Priority: priority}
	if len(extra) > 0 {
		res.Extra = extra[0]
	}
	return res
}

type PermissionDef struct {
	Namespace    string
	Side         PermissionSide
	DisplayName  string
	Action       Action
	Extra        map[string]interface{}
	internalOnly bool
}

func NewPermissionDef(namespace string, action Action, displayName string, side PermissionSide, extra ...map[string]interface{}) *PermissionDef {
	res := &PermissionDef{Namespace: namespace, Action: action, Side: side, DisplayName: displayName}
	if len(extra) > 0 {
		res.Extra = extra[0]
	}
	return res
}

func (d *PermissionDef) AsInternalOnly() *PermissionDef {
	d.internalOnly = true
	return d
}

type PermissionSide int32

const (
	PermissionBothSide       PermissionSide = 0
	PermissionHostSideOnly   PermissionSide = 1
	PermissionTenantSideOnly PermissionSide = 2
)

func AddGroup(group *PermissionDefGroup) *PermissionDefGroup {
	lock.Lock()
	defer lock.Unlock()
	groups = append(groups, group)
	return group
}

func MustFindDef(namespace string, action Action) PermissionDef {
	def, err := FindDef(namespace, action, false)
	if err != nil {
		panic(err)
	}
	return def
}

func FindDef(namespace string, action Action, publicOnly bool) (PermissionDef, error) {
	//TODO cache?
	def, ok := lo.Find(lo.FlatMap(groups, func(t *PermissionDefGroup, _ int) []PermissionDef {
		return t.def
	}), func(t PermissionDef) bool {
		return t.Namespace == namespace && t.Action.GetIdentity() == action.GetIdentity()
	})
	if !ok || (publicOnly && def.internalOnly) {
		return PermissionDef{}, errors.New(400, "PERMISSION_DEF_NOT_FOUND", fmt.Sprintf("action %s in %s not defined", action.GetIdentity(), namespace))
	}
	return def, nil
}

func (g *PermissionDefGroup) AddDef(def *PermissionDef) *PermissionDefGroup {
	if g.Side != PermissionBothSide && g.Side != def.Side {
		panic(fmt.Sprintf("group %s has permission side %v, but try to add permission %s with %v side", g.DisplayName, g.Side, def.DisplayName, def.Side))
	}
	g.def = append(g.def, *def)
	return g
}

func (g *PermissionDefGroup) Walk(isHost bool, publicOnly bool, f func(def PermissionDef)) {
	for _, def := range g.def {
		if publicOnly && def.internalOnly {
			continue
		}
		if (def.Side == PermissionHostSideOnly && !isHost) || (def.Side == PermissionTenantSideOnly && isHost) {
			continue
		}
		f(def)
	}
}

func WalkGroups(isHost bool, f func(group PermissionDefGroup)) {
	var sortedGroup []PermissionDefGroup
	sortedGroup = append(sortedGroup, lo.Map(groups, func(t *PermissionDefGroup, _ int) PermissionDefGroup {
		return *t
	})...)
	sort.SliceStable(sortedGroup, func(i, j int) bool {
		return sortedGroup[i].Priority < sortedGroup[j].Priority
	})
	for _, g := range sortedGroup {
		if (g.Side == PermissionHostSideOnly && !isHost) || (g.Side == PermissionTenantSideOnly && isHost) {
			continue
		}
		f(g)
	}
}