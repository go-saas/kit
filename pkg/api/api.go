package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	shttp "github.com/goxiaoy/go-saas/common/http"
	"github.com/goxiaoy/go-saas/kratos/middleware"
	grpc2 "google.golang.org/grpc"
)

// NewGrpcConn create new grpc client from name
func NewGrpcConn(name string, services *conf.Services, insecure bool,hmtOpt *shttp.WebMultiTenancyOption, opts ...grpc.ClientOption) (grpc2.ClientConnInterface, func()) {
	endpoint, ok := services.Services[name]
	if !ok {
		panic(errors.New(fmt.Sprintf(" %v service not found", name)))
	}
	var conn *grpc2.ClientConn
	var err error
	fOpts := []grpc.ClientOption{
		grpc.WithEndpoint(endpoint.GrpcEndpoint),
		grpc.WithMiddleware(middleware.Client(hmtOpt)),
	}
	if insecure {
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
func NewHttpClient(name string, services *conf.Services,hmtOpt *shttp.WebMultiTenancyOption, opts ...http.ClientOption) (*http.Client, func()) {
	endpoint, ok := services.Services[name]
	if !ok {
		panic(errors.New(fmt.Sprintf(" %v service not found", name)))
	}
	fOpts := []http.ClientOption{
		http.WithEndpoint(endpoint.HttpEndpoint),
		http.WithMiddleware(middleware.Client(hmtOpt)),
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
