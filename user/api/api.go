package api

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	grpc2 "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/authz/authz"
	"github.com/go-saas/kit/pkg/conf"
	kitdi "github.com/go-saas/kit/pkg/di"
	v13 "github.com/go-saas/kit/user/api/account/v1"
	v12 "github.com/go-saas/kit/user/api/auth/v1"
	v15 "github.com/go-saas/kit/user/api/permission/v1"
	v14 "github.com/go-saas/kit/user/api/role/v1"
	v1 "github.com/go-saas/kit/user/api/user/v1"
	_ "github.com/go-saas/kit/user/i18n"
	"github.com/goava/di"
	"google.golang.org/grpc"
)

type GrpcConn grpc.ClientConnInterface
type HttpClient *http.Client

const ServiceName = "user"

func NewGrpcConn(
	client *conf.Client,
	services *conf.Services,
	dis registry.Discovery,
	opt *api.Option,
	tokenMgr api.TokenManager,
	logger log.Logger,
	opts []grpc2.ClientOption,
) (GrpcConn, func()) {
	return api.NewGrpcConn(client, ServiceName, services, dis, opt, tokenMgr, logger, opts)
}

var GrpcProviderSet = kitdi.NewSet(
	NewUserTenantContrib, NewRefreshProvider,
	kitdi.NewProvider(NewRemotePermissionChecker, di.As(new(authz.PermissionChecker)), di.As(new(authz.PermissionManagementService))),
	NewGrpcConn,
	NewUserGrpcClient, NewUserInternalGrpcClient, NewAuthGrpcClient, NewAccountGrpcClient, NewRoleGrpcClient, NewPermissionGrpcClient, NewPermissionInternalGrpcClient)

func NewUserGrpcClient(conn GrpcConn) v1.UserServiceServer {
	return v1.NewUserServiceClientProxy(v1.NewUserServiceClient(conn))
}

func NewUserInternalGrpcClient(conn GrpcConn) v1.UserInternalServiceServer {
	return v1.NewUserInternalServiceClientProxy(v1.NewUserInternalServiceClient(conn))
}

func NewPermissionGrpcClient(conn GrpcConn) v15.PermissionServiceServer {
	return v15.NewPermissionServiceClientProxy(v15.NewPermissionServiceClient(conn))
}

func NewPermissionInternalGrpcClient(conn GrpcConn) v15.PermissionInternalServiceServer {
	return v15.NewPermissionInternalServiceClientProxy(v15.NewPermissionInternalServiceClient(conn))
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
