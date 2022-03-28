package api

import (
	"fmt"

	"github.com/goxiaoy/go-saas-kit/pkg/authz/authz"
)

const (
	ResourceMenu = "sys.menu"
	ResourceDev  = "dev"

	ResourceDevJaeger = "dev.jaeger"
)

func init() {
	//menu management is only allowed in host side
	authz.AddGroup(authz.NewPermissionDefGroup(fmt.Sprintf("%s.permission", ResourceMenu), authz.PermissionHostSideOnly, 0).
		AddDef(authz.NewPermissionDef(ResourceMenu, authz.AnyAction, fmt.Sprintf("%s.any", ResourceMenu), authz.PermissionHostSideOnly).AsInternalOnly()).
		AddDef(authz.NewPermissionDef(ResourceMenu, authz.ReadAction, fmt.Sprintf("%s.read", ResourceMenu), authz.PermissionHostSideOnly)).
		AddDef(authz.NewPermissionDef(ResourceMenu, authz.CreateAction, fmt.Sprintf("%s.create", ResourceMenu), authz.PermissionHostSideOnly)).
		AddDef(authz.NewPermissionDef(ResourceMenu, authz.UpdateAction, fmt.Sprintf("%s.update", ResourceMenu), authz.PermissionHostSideOnly)).
		AddDef(authz.NewPermissionDef(ResourceMenu, authz.DeleteAction, fmt.Sprintf("%s.delete", ResourceMenu), authz.PermissionHostSideOnly)))

	authz.AddGroup(authz.NewPermissionDefGroup(ResourceDev, authz.PermissionHostSideOnly, 0)).
		AddDef(authz.NewPermissionDef(ResourceDev, authz.AnyAction, fmt.Sprintf("%s.any", ResourceDev), authz.PermissionHostSideOnly)).
		AddDef(authz.NewPermissionDef(ResourceDevJaeger, authz.AnyAction, fmt.Sprintf("%s.any", ResourceDevJaeger), authz.PermissionHostSideOnly))
}
