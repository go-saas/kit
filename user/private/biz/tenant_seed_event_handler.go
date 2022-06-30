package biz

import (
	"context"
	"fmt"
	"github.com/go-saas/kit/event"
	v1 "github.com/go-saas/kit/saas/event/v1"
	"github.com/go-saas/kit/user/api"
	"github.com/go-saas/saas/seed"
	"github.com/hibiken/asynq"
	"google.golang.org/protobuf/encoding/protojson"
	"time"
)

type TenantSeedEventHandler event.ConsumerHandler

func NewTenantSeedEventHandler(client *asynq.Client) TenantSeedEventHandler {
	msg := &v1.TenantCreatedEvent{}
	return event.ProtoHandler[*v1.TenantCreatedEvent](msg, event.HandlerFuncOf[*v1.TenantCreatedEvent](func(ctx context.Context, msg *v1.TenantCreatedEvent) error {
		t, err := NewUserMigrationTask(msg)
		if err != nil {
			return err
		}
		_, err = client.EnqueueContext(ctx, t)
		return err
	}))
}

func NewUserMigrationTask(msg *v1.TenantCreatedEvent) (*asynq.Task, error) {
	payload, err := protojson.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(JobTypeUserMigration, payload, asynq.Queue(string(ConnName)), asynq.Retention(time.Hour*24*30)), nil
}

type UserMigrationTaskHandler func(ctx context.Context, t *asynq.Task) error

func NewUserMigrationTaskHandler(seeder seed.Seeder, sender event.Producer) UserMigrationTaskHandler {
	return func(ctx context.Context, t *asynq.Task) error {
		msg := &v1.TenantCreatedEvent{}
		if err := protojson.Unmarshal(t.Payload(), msg); err != nil {
			return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
		}
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
		if err := seeder.Seed(ctx, seed.AddTenant(msg.Id), seed.WithExtra(extra)); err != nil {
			return err
		}
		e := &v1.TenantReadyEvent{
			Id:          msg.Id,
			ServiceName: api.ServiceName,
		}
		ee, _ := event.NewMessageFromProto(e)
		return sender.Send(ctx, ee)
	}
}
