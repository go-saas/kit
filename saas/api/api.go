package api

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	grpc2 "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/conf"
	v1 "github.com/go-saas/kit/saas/api/tenant/v1"
	"github.com/google/wire"
	"google.golang.org/grpc"
)

type GrpcConn grpc.ClientConnInterface
type HttpClient *http.Client

const ServiceName = "saas"

func NewGrpcConn(clientName api.ClientName, services *conf.Services, dis registry.Discovery,
	opt *api.Option, tokenMgr api.TokenManager,
	logger log.Logger,
	opts ...grpc2.ClientOption) (GrpcConn, func()) {
	return api.NewGrpcConn(clientName, ServiceName, services, dis, opt, tokenMgr, logger, opts...)
}

var GrpcProviderSet = wire.NewSet(NewTenantStore, NewGrpcConn, NewTenantGrpcClient, NewTenantInternalGrpcClient)

func NewTenantGrpcClient(conn GrpcConn) v1.TenantServiceServer {
	return v1.NewTenantServiceClientProxy(v1.NewTenantServiceClient(conn))
}

func NewTenantInternalGrpcClient(conn GrpcConn) v1.TenantInternalServiceServer {
	return v1.NewTenantInternalServiceClientProxy(v1.NewTenantInternalServiceClient(conn))
}
