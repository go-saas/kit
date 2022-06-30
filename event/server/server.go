package server

import (
	v12 "github.com/go-saas/kit/event/api/v1"
	"github.com/go-saas/kit/event/service"
	"github.com/google/wire"
)

var EventProviderSet = wire.NewSet(
	service.NewEventService, wire.Bind(new(v12.EventServiceServer), new(*service.EventService)),
)
