package server

import (
	dtmapi "github.com/go-saas/kit/dtm/api"
	v1 "github.com/go-saas/kit/dtm/api/dtm/v1"
	"github.com/go-saas/kit/dtm/data"
	"github.com/go-saas/kit/dtm/service"
	kitdi "github.com/go-saas/kit/pkg/di"
	"github.com/goava/di"
)

var DtmProviderSet = kitdi.NewSet(
	dtmapi.NewInit,
	kitdi.NewProvider(service.NewMsgService, di.As(new(v1.MsgServiceServer))),
	data.NewBarrierMigrator, data.NewStorageMigrator, data.NewMigrator,
	service.NewHttpServerRegister, service.NewGrpcServerRegister,
)
