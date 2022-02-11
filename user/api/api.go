package api

import (
	grpc2 "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/pkg/api"
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

func NewGrpcConn(clientName api.ClientName, services *conf.Services, opt *api.Option, tokenMgr api.TokenManager, opts ...grpc2.ClientOption) (GrpcConn, func()) {
	return api.NewGrpcConn(clientName, ServiceName, services, true, opt, tokenMgr, opts...)
}

func NewHttpClient(clientName api.ClientName, services *conf.Services, opt *api.Option, tokenMgr api.TokenManager, opts ...http.ClientOption) (HttpClient, func()) {
	return api.NewHttpClient(clientName, ServiceName, services, opt, tokenMgr, opts...)
}

var GrpcProviderSet = wire.NewSet(NewGrpcConn, NewUserGrpcClient, NewAuthGrpcClient, NewAccountGrpcClient, NewRoleGrpcClient, NewPermissionGrpcClient)
var HttpProviderSet = wire.NewSet(NewHttpClient, NewUserHttpClient, NewAuthHttpClient, NewAccountHttpClient, NewRoleHttpClient, NewPermissionHttpClient)

func NewUserGrpcClient(conn GrpcConn) v1.UserServiceClient {
	return v1.NewUserServiceClient(conn)
}

func NewPermissionGrpcClient(conn GrpcConn) v15.PermissionServiceClient {
	return v15.NewPermissionServiceClient(conn)
}

func NewPermissionHttpClient(http HttpClient) v15.PermissionServiceHTTPClient {
	return v15.NewPermissionServiceHTTPClient(http)
}

func NewUserHttpClient(http HttpClient) v1.UserServiceHTTPClient {
	return v1.NewUserServiceHTTPClient(http)
}

func NewAuthGrpcClient(conn GrpcConn) v12.AuthClient {
	return v12.NewAuthClient(conn)
}

func NewAuthHttpClient(http HttpClient) v12.AuthHTTPClient {
	return v12.NewAuthHTTPClient(http)
}

func NewAccountGrpcClient(conn GrpcConn) v13.AccountClient {
	return v13.NewAccountClient(conn)
}

func NewAccountHttpClient(http HttpClient) v13.AccountHTTPClient {
	return v13.NewAccountHTTPClient(http)
}

func NewRoleGrpcClient(conn GrpcConn) v14.RoleServiceClient {
	return v14.NewRoleServiceClient(conn)
}

func NewRoleHttpClient(http HttpClient) v14.RoleServiceHTTPClient {
	return v14.NewRoleServiceHTTPClient(http)
}
