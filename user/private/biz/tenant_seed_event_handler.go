package biz

import (
	"context"
	"github.com/goxiaoy/go-saas-kit/pkg/event/event"
	v1 "github.com/goxiaoy/go-saas-kit/saas/event/v1"
	"github.com/goxiaoy/go-saas/seed"
)

type TenantSeedEventHandler event.Handler

func NewTenantSeedEventHandler(seeder seed.Seeder) TenantSeedEventHandler {
	msg := &v1.TenantCreatedEvent{}
	return TenantSeedEventHandler(event.ProtoHandler[*v1.TenantCreatedEvent](msg, func(ctx context.Context, msg *v1.TenantCreatedEvent) error {
		//user seed ignore separate db due to not support
		return seeder.Seed(ctx, seed.NewSeedOption().WithTenantId(msg.Id).WithExtra(map[string]interface{}{SkipMigration: true}))
	}))
}
