package remote

import (
	"github.com/google/wire"
)

var GrpcProviderSet = wire.NewSet(NewRemoteGrpcTenantStore)
