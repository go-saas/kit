package server

import (
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-saas/kit/pkg/job"
	sbiz "github.com/go-saas/kit/saas/private/biz"
	sysbiz "github.com/go-saas/kit/sys/private/biz"
	ubiz "github.com/go-saas/kit/user/private/biz"
	"github.com/hibiken/asynq"
)

func NewJobServer(opt asynq.RedisConnOpt, log klog.Logger, handler ubiz.UserMigrationTaskHandler) *job.Server {
	srv := job.NewServer(opt, job.WithQueues(map[string]int{
		string(ubiz.ConnName):   1,
		string(sbiz.ConnName):   1,
		string(sysbiz.ConnName): 1,
	}))
	srv.Use(job.TracingServer(), job.Logging(log), job.Stack())
	srv.HandleFunc(ubiz.JobTypeUserMigration, handler)
	return srv
}
