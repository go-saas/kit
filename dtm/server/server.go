package server

import (
	dtmapi "github.com/go-saas/kit/dtm/api"
	v1 "github.com/go-saas/kit/dtm/api/dtm/v1"
	"github.com/go-saas/kit/dtm/data"
	"github.com/go-saas/kit/dtm/service"
	"github.com/google/wire"
)

var DtmProviderSet = wire.NewSet(
	dtmapi.NewInit,
	service.NewMsgService, wire.Bind(new(v1.MsgServiceServer), new(*service.MsgService)),
	data.NewMigrator,
	service.NewHttpServerRegister, service.NewGrpcServerRegister,
)
