package server

import (
	"github.com/goxiaoy/go-saas-kit/pkg/job"
	sbiz "github.com/goxiaoy/go-saas-kit/saas/private/biz"
	sysbiz "github.com/goxiaoy/go-saas-kit/sys/private/biz"
	ubiz "github.com/goxiaoy/go-saas-kit/user/private/biz"
	"github.com/hibiken/asynq"
)

func NewJobServer(opt asynq.RedisConnOpt, handler ubiz.UserMigrationTaskHandler) *job.Server {
	srv := job.NewServer(opt, job.WithQueues(map[string]int{
		string(ubiz.ConnName):   1,
		string(sbiz.ConnName):   1,
		string(sysbiz.ConnName): 1,
	}))
	srv.HandleFunc(ubiz.JobTypeUserMigration, handler)
	return srv
}
