package keto

import (
	grpc2 "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/pkg/api"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authorization"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	acl "github.com/ory/keto/proto/ory/keto/acl/v1alpha1"
	"google.golang.org/grpc"
)

const ServiceName = "keto"

type GrpcConn grpc.ClientConnInterface

func NewGrpcConn(clientName api.ClientName, services *conf.Services, opt *api.Option, tokenMgr api.TokenManager, opts ...grpc2.ClientOption) (GrpcConn, func()) {
	return api.NewGrpcConn(clientName, ServiceName, services, true, opt, tokenMgr, opts...)
}

func NewCheckServiceClient(conn GrpcConn) acl.CheckServiceClient {
	return NewCheckServiceClient(conn)
}

var ProviderSet = wire.NewSet(NewGrpcConn, NewCheckServiceClient, NewPermissionChecker, wire.Bind(new(authorization.PermissionChecker), new(*PermissionChecker)))
