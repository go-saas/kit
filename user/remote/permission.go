package remote

import (
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authz"
)

var GrpcProviderSet = wire.NewSet(NewRemotePermissionChecker,
	wire.Bind(new(authz.PermissionChecker), new(*PermissionChecker)),
	wire.Bind(new(authz.PermissionManagementService), new(*PermissionChecker)))
