package api

import (
	"fmt"

	"github.com/goxiaoy/go-saas-kit/pkg/authz/authz"
)

const (
	ResourcePermission         = "permission"
	ResourcePermissionInternal = ResourcePermission + ".internal"

	ResourceUser       = "user.user"
	ResourceUserTenant = "user.user_tenant"
	ResourceRole       = "user.role"
)

func init() {
	authz.AddGroup(authz.NewPermissionDefGroup(fmt.Sprintf("%s.permission", ResourcePermission), authz.PermissionBothSide, 0).
		AddDef(authz.NewPermissionDef(ResourcePermissionInternal, authz.AnyAction, fmt.Sprintf("%s.internal", ResourcePermission), authz.PermissionHostSideOnly).AsInternalOnly()).
		AddDef(authz.NewPermissionDef(ResourcePermission, authz.ReadAction, fmt.Sprintf("%s.read", ResourcePermission), authz.PermissionBothSide)).
		AddDef(authz.NewPermissionDef(ResourcePermission, authz.WriteAction, fmt.Sprintf("%s.write", ResourcePermission), authz.PermissionBothSide)))

	authz.AddGroup(authz.NewPermissionDefGroup(fmt.Sprintf("%s.permission", ResourceUser), authz.PermissionBothSide, 0).
		AddDef(authz.NewPermissionDef(ResourceUserTenant, authz.AnyAction, fmt.Sprintf("%s.any", ResourceUserTenant), authz.PermissionHostSideOnly).AsInternalOnly()).
		AddDef(authz.NewPermissionDef(ResourceUser, authz.AnyAction, fmt.Sprintf("%s.any", ResourceUser), authz.PermissionBothSide).AsInternalOnly()).
		AddDef(authz.NewPermissionDef(ResourceUser, authz.ReadAction, fmt.Sprintf("%s.read", ResourceUser), authz.PermissionBothSide)).
		AddDef(authz.NewPermissionDef(ResourceUser, authz.CreateAction, fmt.Sprintf("%s.create", ResourceUser), authz.PermissionBothSide)).
		AddDef(authz.NewPermissionDef(ResourceUser, authz.UpdateAction, fmt.Sprintf("%s.update", ResourceUser), authz.PermissionBothSide)).
		AddDef(authz.NewPermissionDef(ResourceUser, authz.DeleteAction, fmt.Sprintf("%s.delete", ResourceUser), authz.PermissionBothSide)))

	authz.AddGroup(authz.NewPermissionDefGroup(fmt.Sprintf("%s.permission", ResourceRole), authz.PermissionBothSide, 0).
		AddDef(authz.NewPermissionDef(ResourceRole, authz.AnyAction, fmt.Sprintf("%s.any", ResourceRole), authz.PermissionBothSide).AsInternalOnly()).
		AddDef(authz.NewPermissionDef(ResourceRole, authz.ReadAction, fmt.Sprintf("%s.read", ResourceRole), authz.PermissionBothSide)).
		AddDef(authz.NewPermissionDef(ResourceRole, authz.CreateAction, fmt.Sprintf("%s.create", ResourceRole), authz.PermissionBothSide)).
		AddDef(authz.NewPermissionDef(ResourceRole, authz.UpdateAction, fmt.Sprintf("%s.update", ResourceRole), authz.PermissionBothSide)).
		AddDef(authz.NewPermissionDef(ResourceRole, authz.DeleteAction, fmt.Sprintf("%s.delete", ResourceRole), authz.PermissionBothSide)))
}
