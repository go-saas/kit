package api

import (
	"fmt"

	"github.com/goxiaoy/go-saas-kit/pkg/authz/authz"
)

const (
	ResourceTenant = "saas.tenant"
)

func init() {
	//tenant management is only allowed in host side
	authz.AddGroup(authz.NewPermissionDefGroup(fmt.Sprintf("%s.permission", ResourceTenant), authz.PermissionHostSideOnly, 0).
		AddDef(authz.NewPermissionDef(ResourceTenant, authz.AnyAction, fmt.Sprintf("%s.any", ResourceTenant), authz.PermissionHostSideOnly).AsInternalOnly()).
		AddDef(authz.NewPermissionDef(ResourceTenant, authz.ReadAction, fmt.Sprintf("%s.read", ResourceTenant), authz.PermissionHostSideOnly)).
		AddDef(authz.NewPermissionDef(ResourceTenant, authz.CreateAction, fmt.Sprintf("%s.create", ResourceTenant), authz.PermissionHostSideOnly)).
		AddDef(authz.NewPermissionDef(ResourceTenant, authz.UpdateAction, fmt.Sprintf("%s.update", ResourceTenant), authz.PermissionHostSideOnly)).
		AddDef(authz.NewPermissionDef(ResourceTenant, authz.DeleteAction, fmt.Sprintf("%s.delete", ResourceTenant), authz.PermissionHostSideOnly)))

}
