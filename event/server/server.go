package server

import (
	v12 "github.com/go-saas/kit/event/api/v1"
	"github.com/go-saas/kit/event/service"
	kitdi "github.com/go-saas/kit/pkg/di"
	"github.com/goava/di"
)

var EventProviderSet = kitdi.NewSet(
	kitdi.NewProvider(service.NewEventService, di.As(new(v12.EventServiceServer))),
	service.NewGrpcServerRegister, service.NewHttpServerRegister,
)
