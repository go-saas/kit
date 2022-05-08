package api

import (
	"github.com/go-kratos/kratos/v2/log"
	grpc2 "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
	v1 "github.com/goxiaoy/go-saas-kit/payment/api/post/v1"
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

var GrpcProviderSet = wire.NewSet(NewGrpcConn, NewPostGrpcClient)
var HttpProviderSet = wire.NewSet(NewHttpClient, NewPostHttpClient)

func NewPostGrpcClient(conn GrpcConn) v1.PostServiceClient {
	return v1.NewPostServiceClient(conn)
}

func NewPostHttpClient(http HttpClient) v1.PostServiceHTTPClient {
	return v1.NewPostServiceHTTPClient(http)
}
