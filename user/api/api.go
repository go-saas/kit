package api

import (
	"github.com/go-kratos/kratos/v2/log"
	grpc2 "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/pkg/api"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authz"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	v13 "github.com/goxiaoy/go-saas-kit/user/api/account/v1"
	v12 "github.com/goxiaoy/go-saas-kit/user/api/auth/v1"
	v15 "github.com/goxiaoy/go-saas-kit/user/api/permission/v1"
	v14 "github.com/goxiaoy/go-saas-kit/user/api/role/v1"
	v1 "github.com/goxiaoy/go-saas-kit/user/api/user/v1"
	"google.golang.org/grpc"
)

type GrpcConn grpc.ClientConnInterface
type HttpClient *http.Client

const ServiceName = "user"

func NewGrpcConn(clientName api.ClientName, services *conf.Services, opt *api.Option, tokenMgr api.TokenManager, logger log.Logger, opts ...grpc2.ClientOption) (GrpcConn, func()) {
	return api.NewGrpcConn(clientName, ServiceName, services, opt, tokenMgr, logger, opts...)
}

var GrpcProviderSet = wire.NewSet(
	NewUserTenantContrib, NewRefreshProvider,
	NewRemotePermissionChecker,
	wire.Bind(new(authz.PermissionChecker), new(*PermissionChecker)),
	wire.Bind(new(authz.PermissionManagementService), new(*PermissionChecker)),
	NewGrpcConn,
	NewUserGrpcClient, NewAuthGrpcClient, NewAccountGrpcClient, NewRoleGrpcClient, NewPermissionGrpcClient)

func NewUserGrpcClient(conn GrpcConn) v1.UserServiceServer {
	return v1.NewUserServiceClientProxy(v1.NewUserServiceClient(conn))
}

func NewPermissionGrpcClient(conn GrpcConn) v15.PermissionServiceServer {
	return v15.NewPermissionServiceClientProxy(v15.NewPermissionServiceClient(conn))
}

func NewAuthGrpcClient(conn GrpcConn) v12.AuthServer {
	return v12.NewAuthClientProxy(v12.NewAuthClient(conn))
}

func NewAccountGrpcClient(conn GrpcConn) v13.AccountServer {
	return v13.NewAccountClientProxy(v13.NewAccountClient(conn))
}

func NewRoleGrpcClient(conn GrpcConn) v14.RoleServiceServer {
	return v14.NewRoleServiceClientProxy(v14.NewRoleServiceClient(conn))
}
