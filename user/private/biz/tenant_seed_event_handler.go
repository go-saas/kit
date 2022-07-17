package biz

import (
	"context"
	"fmt"
	"github.com/go-saas/kit/event"
	"github.com/go-saas/kit/pkg/job"
	v1 "github.com/go-saas/kit/saas/event/v1"
	"github.com/go-saas/kit/user/api"
	"github.com/go-saas/saas/seed"
	"github.com/hibiken/asynq"
	"google.golang.org/protobuf/encoding/protojson"
	"time"
)

const (
	JobTypeUserMigration = string(ConnName) + ":" + "migration"
)

func NewTenantSeedEventHandler(client *asynq.Client) event.ConsumerHandler {
	msg := &v1.TenantCreatedEvent{}
	return event.ProtoHandler[*v1.TenantCreatedEvent](msg, event.HandlerFuncOf[*v1.TenantCreatedEvent](func(ctx context.Context, msg *v1.TenantCreatedEvent) error {
		t, err := NewUserMigrationTask(NewUserMigrationTaskFromTenantEvent(msg))
		if err != nil {
			return err
		}
		_, err = client.EnqueueContext(ctx, t)
		return err
	}))
}

func NewUserMigrationTask(msg *UserMigrationTask) (*asynq.Task, error) {
	payload, err := protojson.Marshal(msg)
	if err != nil {
		return nil, err
	}
	// delay second in case saas local transaction not committed
	return asynq.NewTask(JobTypeUserMigration, payload, asynq.ProcessIn(time.Second), asynq.Queue(string(ConnName)), asynq.Retention(time.Hour*24*30)), nil
}

func NewUserMigrationTaskHandler(seeder seed.Seeder, producer event.Producer) *job.Handler {
	return job.NewHandlerFunc(JobTypeUserMigration, func(ctx context.Context, t *asynq.Task) error {
		msg := &UserMigrationTask{}
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
		if len(msg.AdminUserId) > 0 {
			extra[AdminUserId] = msg.AdminUserId
		}
		if err := seeder.Seed(ctx, seed.AddTenant(msg.Id), seed.WithExtra(extra)); err != nil {
			return err
		}
		e := &v1.TenantReadyEvent{
			Id:          msg.Id,
			ServiceName: api.ServiceName,
		}
		ee, _ := event.NewMessageFromProto(e)
		return producer.Send(ctx, ee)
	})
}
