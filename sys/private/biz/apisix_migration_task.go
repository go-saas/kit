package biz

import (
	"context"
	"github.com/go-saas/kit/pkg/job"
	"github.com/hibiken/asynq"
	"time"
)

const (
	JobTypeApisixMigration = string(ConnName) + ":" + "apisix" + ":" + "migration"
)

func NewApisixMigrationTask() *asynq.Task {
	return asynq.NewTask(JobTypeApisixMigration, nil, asynq.ProcessIn(time.Second), asynq.Queue(string(ConnName)), asynq.Retention(time.Hour*24*30))
}

func NewApisixMigrationTaskHandler(seeder *ApisixSeed) *job.Handler {
	return job.NewHandlerFunc(JobTypeApisixMigration, func(ctx context.Context, t *asynq.Task) error {
		return seeder.Do(ctx)
	})
}
