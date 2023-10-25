package api

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	grpc2 "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/conf"
	kitdi "github.com/go-saas/kit/pkg/di"
	v1 "github.com/go-saas/kit/product/api/product/v1"
	_ "github.com/go-saas/kit/product/i18n"
	"google.golang.org/grpc"
)

type GrpcConn grpc.ClientConnInterface
type HttpClient *http.Client

const (
	ServiceName = "product"
)

func NewGrpcConn(client *conf.Client, services *conf.Services, dis registry.Discovery, opt *api.Option, tokenMgr api.TokenManager, logger log.Logger, opts []grpc2.ClientOption) (GrpcConn, func()) {
	return api.NewGrpcConn(client, ServiceName, services, dis, opt, tokenMgr, logger, opts)
}

var GrpcProviderSet = kitdi.NewSet(NewGrpcConn, NewProductGrpcClient, NewProductInternalGrpcClient)

func NewProductGrpcClient(conn GrpcConn) v1.ProductServiceServer {
	return v1.NewProductServiceClientProxy(v1.NewProductServiceClient(conn))
}
func NewProductInternalGrpcClient(conn GrpcConn) v1.ProductInternalServiceServer {
	return v1.NewProductInternalServiceClientProxy(v1.NewProductInternalServiceClient(conn))
}
