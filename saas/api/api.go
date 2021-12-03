package api

import (
	grpc2 "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/pkg/api"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	v1 "github.com/goxiaoy/go-saas-kit/saas/api/tenant/v1"
	shttp "github.com/goxiaoy/go-saas/common/http"

	"google.golang.org/grpc"
)

// depend on user to prevent circle dependency
import _ "github.com/goxiaoy/go-saas-kit/user/api"

type GrpcConn grpc.ClientConnInterface
type HttpClient *http.Client

const ServiceName = "saas"

func NewGrpcConn(services *conf.Services,hmtOpt *shttp.WebMultiTenancyOption, opts ...grpc2.ClientOption) (GrpcConn, func()) {
	return api.NewGrpcConn(ServiceName, services, true,hmtOpt, opts...)
}

func NewHttpClient(services *conf.Services,hmtOpt *shttp.WebMultiTenancyOption, opts ...http.ClientOption) (HttpClient, func()) {
	return api.NewHttpClient(ServiceName, services,hmtOpt, opts...)
}

var GrpcProviderSet = wire.NewSet(NewGrpcConn, NewTenantGrpcClient, NewRemoteGrpcTenantStore)
var HttpProviderSet = wire.NewSet(NewHttpClient, NewTenantHttpClient, NewRemoteHttpTenantStore)

func NewTenantGrpcClient(conn GrpcConn) v1.TenantServiceClient {
	return v1.NewTenantServiceClient(conn)
}

func NewTenantHttpClient(http HttpClient) v1.TenantServiceHTTPClient {
	return v1.NewTenantServiceHTTPClient(http)
}
