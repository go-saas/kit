package common

import (
	"github.com/go-saas/kit/pkg/conf"
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
