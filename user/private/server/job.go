package server

import (
	"github.com/goxiaoy/go-saas-kit/pkg/job"
	"github.com/goxiaoy/go-saas-kit/user/private/biz"
	"github.com/hibiken/asynq"
)

func NewJobServer(opt asynq.RedisConnOpt, handler biz.UserMigrationTaskHandler) *job.Server {
	srv := job.NewServer(opt, job.WithQueues(map[string]int{
		string(biz.ConnName): 1,
	}))
	srv.HandleFunc(biz.JobTypeUserMigration, handler)
	return srv
}
