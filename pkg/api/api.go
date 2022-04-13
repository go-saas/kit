package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	grpcx "google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"time"
)

type ClientName string

var (
	defaultClientCfg = &conf.Client{Timeout: durationpb.New(5 * time.Second)}
)

// NewGrpcConn create new grpc client from name
func NewGrpcConn(
	clientName ClientName,
	serviceName string,
	services *conf.Services,
	opt *Option,
	tokenMgr TokenManager,
	logger log.Logger,
	opts ...grpc.ClientOption,
) (grpcx.ClientConnInterface, func()) {
	endpoint, ok := services.Services[serviceName]
	if !ok {
		panic(errors.New(fmt.Sprintf(" %v service not found", serviceName)))
	}

	var clientCfg = proto.Clone(defaultClientCfg).(*conf.Client)

	if c, ok := services.Clients[string(clientName)]; !ok {
		panic(errors.New(fmt.Sprintf(" %v client not found", clientName)))
	} else {
		proto.Merge(clientCfg, c)
	}
	var conn *grpcx.ClientConn
	var err error
	fOpts := []grpc.ClientOption{
		grpc.WithEndpoint(endpoint.GrpcEndpoint),
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
func NewHttpClient(clientName ClientName, serviceName string,
	services *conf.Services, opt *Option, tokenMgr TokenManager,
	logger log.Logger,
	opts ...http.ClientOption) (*http.Client, func()) {
	endpoint, ok := services.Services[serviceName]
	if !ok {
		panic(errors.New(fmt.Sprintf(" %v service not found", serviceName)))
	}

	clientCfg, ok := services.Clients[string(clientName)]
	if !ok {
		panic(errors.New(fmt.Sprintf(" %v client not found", clientName)))
	}

	fOpts := []http.ClientOption{
		http.WithEndpoint(endpoint.HttpEndpoint),
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

var DefaultProviderSet = wire.NewSet(
	NewDefaultOption,
	NewInMemoryTokenManager,
	wire.Bind(new(TokenManager), new(*InMemoryTokenManager)),
	NewClientTrustedContextValidator,
)
