package server

import (
	v1 "github.com/go-saas/kit/dtm/api/dtm/v1"
	"github.com/go-saas/kit/dtm/data"
	"github.com/go-saas/kit/dtm/service"
	"github.com/google/wire"
)

var DtmProviderSet = wire.NewSet(
	service.NewMsgService, wire.Bind(new(v1.MsgServiceServer), new(*service.MsgServiceService)),
	data.NewMigrator,
	service.NewHttpServerRegister, service.NewGrpcServerRegister,
)
