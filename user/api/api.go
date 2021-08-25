package api

import (
	grpc2 "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/pkg/api"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	v13 "github.com/goxiaoy/go-saas-kit/user/api/account/v1"
	v12 "github.com/goxiaoy/go-saas-kit/user/api/auth/v1"
	v1 "github.com/goxiaoy/go-saas-kit/user/api/user/v1"
	"google.golang.org/grpc"
)

type GrpcConn grpc.ClientConnInterface
type HttpClient *http.Client

const ServiceName = "user"

func NewGrpcConn(services *conf.Services, opts ...grpc2.ClientOption) (GrpcConn, func()) {
	return  api.NewGrpcConn(ServiceName, services, true, opts...)
}

func NewHttpClient(services *conf.Services, opts ...http.ClientOption) (HttpClient, func()) {
	return api.NewHttpClient(ServiceName, services, opts...)
}

var GrpcProviderSet = wire.NewSet(NewGrpcConn, NewUserGrpcClient, NewAuthGrpcClient, NewAccountGrpcClient)
var HttpProviderSet = wire.NewSet(NewHttpClient, NewUserHttpClient, NewAuthHttpClient, NewAccountHttpClient)

func NewUserGrpcClient(conn GrpcConn) v1.UserServiceClient {
	return v1.NewUserServiceClient(conn)
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
