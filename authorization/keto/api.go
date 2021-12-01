package keto

import (
	grpc2 "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/authorization/common"
	"github.com/goxiaoy/go-saas-kit/pkg/api"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	shttp "github.com/goxiaoy/go-saas/common/http"
	acl "github.com/ory/keto/proto/ory/keto/acl/v1alpha1"
	"google.golang.org/grpc"
)

const ServiceName = "keto"

type GrpcConn grpc.ClientConnInterface

func NewGrpcConn(services *conf.Services, hmtOpt *shttp.WebMultiTenancyOption, opts ...grpc2.ClientOption) (GrpcConn, func()) {
	return api.NewGrpcConn(ServiceName, services, true, hmtOpt, opts...)
}

func NewCheckServiceClient(conn GrpcConn) acl.CheckServiceClient {
	return NewCheckServiceClient(conn)
}

var ProviderSet = wire.NewSet(NewGrpcConn, NewCheckServiceClient, NewPermissionChecker, wire.Bind(new(common.PermissionChecker), new(*PermissionChecker)))
