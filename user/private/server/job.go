package server

import (
	"github.com/goxiaoy/go-saas-kit/pkg/job"
	"github.com/goxiaoy/go-saas-kit/user/private/biz"
)

func NewJobServer(s job.LazyAsynqServer, handler biz.UserMigrationTaskHandler) *job.Server {
	srv := job.NewServer(s)
	srv.HandleFunc(biz.JobTypeUserMigration, handler)
	return srv
}
