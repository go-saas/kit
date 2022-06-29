package job

import (
	"github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"
	"net/http"
)

func NewUi(root string, opt asynq.RedisConnOpt) http.Handler {
	h := asynqmon.New(asynqmon.Options{
		RootPath:     root, // RootPath specifies the root for asynqmon app
		RedisConnOpt: opt,
	})
	return h
}
