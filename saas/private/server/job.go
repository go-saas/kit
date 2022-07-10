package server

import (
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-saas/kit/pkg/job"
	"github.com/go-saas/kit/saas/private/biz"
	"github.com/hibiken/asynq"
)

func NewJobServer(opt asynq.RedisConnOpt, log klog.Logger) *job.Server {
	srv := job.NewServer(opt, job.WithQueues(map[string]int{
		string(biz.ConnName): 1,
	}))
	srv.Use(job.TracingServer(), job.Logging(log))
	return srv
}
