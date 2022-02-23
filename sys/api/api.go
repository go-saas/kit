package api

import (
	grpc2 "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/pkg/api"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	v1 "github.com/goxiaoy/go-saas-kit/sys/api/menu/v1"
	"google.golang.org/grpc"
)

type GrpcConn grpc.ClientConnInterface
type HttpClient *http.Client

const ServiceName = "sys"

func NewGrpcConn(clientName api.ClientName, services *conf.Services, opt *api.Option, tokenMgr api.TokenManager, opts ...grpc2.ClientOption) (GrpcConn, func()) {
	return api.NewGrpcConn(clientName, ServiceName, services, true, opt, tokenMgr, opts...)
}

func NewHttpClient(clientName api.ClientName, services *conf.Services, opt *api.Option, tokenMgr api.TokenManager, opts ...http.ClientOption) (HttpClient, func()) {
	return api.NewHttpClient(clientName, ServiceName, services, opt, tokenMgr, opts...)
}

var GrpcProviderSet = wire.NewSet(NewGrpcConn, NewMenuGrpcClient)
var HttpProviderSet = wire.NewSet(NewHttpClient, NewMenuServiceHttpClient)

func NewMenuGrpcClient(conn GrpcConn) v1.MenuServiceClient {
	return v1.NewMenuServiceClient(conn)
}

func NewMenuServiceHttpClient(http HttpClient) v1.MenuServiceHTTPClient {
	return v1.NewMenuServiceHTTPClient(http)
}
