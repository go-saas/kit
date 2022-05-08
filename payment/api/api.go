package api

import (
	"github.com/go-kratos/kratos/v2/log"
	grpc2 "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
	v1 "github.com/goxiaoy/go-saas-kit/payment/api/order/v1"
	"github.com/goxiaoy/go-saas-kit/pkg/api"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	"google.golang.org/grpc"
)

type GrpcConn grpc.ClientConnInterface
type HttpClient *http.Client

const ServiceName = "payment"

func NewGrpcConn(clientName api.ClientName, services *conf.Services, opt *api.Option, tokenMgr api.TokenManager, logger log.Logger, opts ...grpc2.ClientOption) (GrpcConn, func()) {
	return api.NewGrpcConn(clientName, ServiceName, services, opt, tokenMgr, logger, opts...)
}

func NewHttpClient(clientName api.ClientName, services *conf.Services, opt *api.Option, tokenMgr api.TokenManager, logger log.Logger, opts ...http.ClientOption) (HttpClient, func()) {
	return api.NewHttpClient(clientName, ServiceName, services, opt, tokenMgr, logger, opts...)
}

var GrpcProviderSet = wire.NewSet(NewGrpcConn, NewOrderGrpcClient)
var HttpProviderSet = wire.NewSet(NewHttpClient, NewOrderHttpClient)

func NewOrderGrpcClient(conn GrpcConn) v1.OrderServiceClient {
	return v1.NewOrderServiceClient(conn)
}

func NewOrderHttpClient(http HttpClient) v1.OrderServiceHTTPClient {
	return v1.NewOrderServiceHTTPClient(http)
}
