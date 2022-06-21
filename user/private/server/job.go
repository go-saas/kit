package server

import (
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/goxiaoy/go-saas-kit/pkg/job"
	"github.com/goxiaoy/go-saas-kit/user/private/biz"
	"github.com/hibiken/asynq"
)

func NewJobServer(opt asynq.RedisConnOpt, log klog.Logger, handler biz.UserMigrationTaskHandler) *job.Server {
	srv := job.NewServer(opt, job.WithQueues(map[string]int{
		string(biz.ConnName): 1,
	}))
	srv.Use(job.TracingServer(), job.Logging(log), job.Stack())
	srv.HandleFunc(biz.JobTypeUserMigration, handler)
	return srv
}
