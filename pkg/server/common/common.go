package common

import (
	"github.com/go-saas/kit/pkg/conf"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"time"
)

const (
	DefaultSrvName = "default"
)

var (
	DefaultServerConfig = &conf.Server{
		Http: &conf.Server_HTTP{
			Addr:    ":9080",
			Timeout: durationpb.New(5 * time.Second),
		},
		Grpc: &conf.Server_GRPC{
			Addr:    ":9081",
			Timeout: durationpb.New(5 * time.Second),
		},
	}
)

func GetConf(services *conf.Services, name string) *conf.Server {
	//default config
	server := proto.Clone(DefaultServerConfig).(*conf.Server)
	if def, ok := services.Servers[DefaultSrvName]; ok {
		//merge default config
		proto.Merge(server, def)
	}
	if s, ok := services.Servers[name]; ok {
		//merge service config
		proto.Merge(server, s)
	}
	return server
}
