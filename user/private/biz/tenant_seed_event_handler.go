package biz

import (
	"context"
	"github.com/goxiaoy/go-saas-kit/pkg/event/event"
	v1 "github.com/goxiaoy/go-saas-kit/saas/event/v1"
	"github.com/goxiaoy/go-saas-kit/user/api"
	"github.com/goxiaoy/go-saas/seed"
)

type TenantSeedEventHandler event.Handler

func NewTenantSeedEventHandler(seeder seed.Seeder, sender event.Sender) TenantSeedEventHandler {
	msg := &v1.TenantCreatedEvent{}
	return TenantSeedEventHandler(event.ProtoHandler[*v1.TenantCreatedEvent](msg, func(ctx context.Context, msg *v1.TenantCreatedEvent) error {
		//user seed ignore separate db due to not support
		extra := map[string]interface{}{}
		if len(msg.AdminEmail) > 0 {
			extra[AdminEmailKey] = msg.AdminEmail
		}
		if len(msg.AdminUsername) > 0 {
			extra[AdminUsernameKey] = msg.AdminUsername
		}
		if len(msg.AdminPassword) > 0 {
			extra[AdminPasswordKey] = msg.AdminPassword
		}
		if err := seeder.Seed(ctx, seed.NewSeedOption().WithTenantId(msg.Id).WithExtra(extra)); err != nil {
			return nil
		}
		e := &v1.TenantReadyEvent{
			Id:          msg.Id,
			ServiceName: api.ServiceName,
		}
		ee, _ := event.NewMessageFromProto(e)
		return sender.Send(ctx, ee)
	}))
}
