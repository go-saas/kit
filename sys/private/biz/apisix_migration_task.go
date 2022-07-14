package biz

import (
	"context"
	"github.com/hibiken/asynq"
	"time"
)

const (
	JobTypeApisixMigration = string(ConnName) + ":" + "apisix" + ":" + "migration"
)

func NewApisixMigrationTask() *asynq.Task {
	return asynq.NewTask(JobTypeApisixMigration, nil, asynq.ProcessIn(time.Second), asynq.Queue(string(ConnName)), asynq.Retention(time.Hour*24*30))
}

type ApisixMigrationTaskHandler func(ctx context.Context, t *asynq.Task) error

func NewApisixMigrationTaskHandler(seeder *ApisixSeed) ApisixMigrationTaskHandler {
	return func(ctx context.Context, t *asynq.Task) error {
		return seeder.Do()
	}
}
