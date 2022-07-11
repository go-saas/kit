package api

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-saas/kit/pkg/conf"
	kregistry "github.com/go-saas/kit/pkg/registry"
	"github.com/google/wire"
	grpcx "google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"time"
)

type ClientName string

var (
	defaultClientCfg = &conf.Client{Timeout: durationpb.New(5 * time.Second)}
)

func NewClientCfg(clientName ClientName, services *conf.Services) *conf.Client {
	var clientCfg = proto.Clone(defaultClientCfg).(*conf.Client)
	if c, ok := services.Clients[string(clientName)]; ok {
		proto.Merge(clientCfg, c)
	}
	if len(clientCfg.ClientId) == 0 {
		clientCfg.ClientId = string(clientName)
	}
	return clientCfg
}

// NewGrpcConn create new grpc client from name
func NewGrpcConn(
	clientCfg *conf.Client,
	serviceName string,
	services *conf.Services,
	dis registry.Discovery,
	opt *Option,
	tokenMgr TokenManager,
	logger log.Logger,
	opts ...grpc.ClientOption,
) (grpcx.ClientConnInterface, func()) {

	var conn *grpcx.ClientConn
	var err error

	var name = serviceName
	serviceCfg := services.GetServiceMergedDefault(serviceName)
	if serviceCfg != nil && len(serviceCfg.Redirect) > 0 {
		name = serviceCfg.Redirect
	}
	fOpts := []grpc.ClientOption{
		grpc.WithEndpoint("discovery:///" + name),
		grpc.WithDiscovery(dis),
		grpc.WithMiddleware(
			recovery.Recovery(),
			tracing.Client(),
			ClientPropagation(clientCfg, opt, tokenMgr, logger),
		),
		// for tracing remote ip recording
		grpc.WithOptions(grpcx.WithStatsHandler(&tracing.ClientHandler{})),
	}
	if clientCfg.Timeout != nil {
		fOpts = append(fOpts, grpc.WithTimeout(clientCfg.Timeout.AsDuration()))
	}
	if opt.Insecure {
		fOpts = append(fOpts, opts...)
		conn, err = grpc.DialInsecure(context.Background(), fOpts...)
	} else {
		fOpts = append(fOpts, opts...)
		conn, err = grpc.Dial(context.Background(), fOpts...)
	}
	if err != nil {
		panic(err)
	}
	return conn, func() {
		conn.Close()
	}
}

// NewHttpClient create new http client from name
func NewHttpClient(
	clientCfg *conf.Client,
	serviceName string,
	services *conf.Services,
	dis registry.Discovery,
	opt *Option,
	tokenMgr TokenManager,
	logger log.Logger,
	opts ...http.ClientOption,
) (*http.Client, func()) {

	var name = serviceName
	serviceCfg := services.GetServiceMergedDefault(serviceName)
	if serviceCfg != nil && len(serviceCfg.Redirect) > 0 {
		name = serviceCfg.Redirect
	}

	fOpts := []http.ClientOption{
		http.WithEndpoint("discovery:///" + name),
		http.WithDiscovery(dis),
		http.WithMiddleware(
			recovery.Recovery(),
			tracing.Client(),
			ClientPropagation(clientCfg, opt, tokenMgr, logger),
		),
	}
	if clientCfg.Timeout != nil {
		fOpts = append(fOpts, http.WithTimeout(clientCfg.Timeout.AsDuration()))
	}
	fOpts = append(fOpts, opts...)
	fOpts = append(fOpts, http.WithBlock())
	conn, err := http.NewClient(context.Background(), fOpts...)
	if err != nil {
		panic(err)
	}
	return conn, func() {
		conn.Close()
	}
}

func NewDiscovery(services *conf.Services) (registry.Discovery, error) {
	_, r, err := kregistry.NewRegister(services.Registry)
	return r, err
}

var DefaultProviderSet = wire.NewSet(
	NewClientCfg,
	NewDefaultOption,
	NewInMemoryTokenManager, wire.Bind(new(TokenManager), new(*InMemoryTokenManager)),
	NewClientTrustedContextValidator,
	NewDiscovery,
)
