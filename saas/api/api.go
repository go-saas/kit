package api

import (
	grpc2 "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/pkg/api"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	v1 "github.com/goxiaoy/go-saas-kit/saas/api/tenant/v1"

	"google.golang.org/grpc"
)

type GrpcConn grpc.ClientConnInterface
type HttpClient *http.Client

const ServiceName = "saas"

func NewGrpcConn(services *conf.Services, opts ...grpc2.ClientOption) (GrpcConn, func()) {
	return api.NewGrpcConn(ServiceName, services, true, opts...)
}

func NewHttpClient(services *conf.Services, opts ...http.ClientOption) (HttpClient, func()) {
	return api.NewHttpClient(ServiceName, services, opts...)
}

var GrpcProviderSet = wire.NewSet(NewGrpcConn, NewTenantGrpcClient,NewRemoteGrpcTenantStore)
var HttpProviderSet = wire.NewSet(NewHttpClient, NewTenantHttpClient,NewRemoteHttpTenantStore)

func NewTenantGrpcClient(conn GrpcConn) v1.TenantServiceClient {
	return v1.NewTenantServiceClient(conn)
}

func NewTenantHttpClient(http HttpClient) v1.TenantServiceHTTPClient {
	return v1.NewTenantServiceHTTPClient(http)
}
