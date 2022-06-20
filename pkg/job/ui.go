package job

import (
	kerrors "github.com/go-kratos/kratos/v2/errors"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
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

func abortWithError(err error, w http.ResponseWriter) {
	//use error codec
	fr := kerrors.FromError(err)
	w.WriteHeader(int(fr.Code))
	khttp.DefaultErrorEncoder(w, &http.Request{}, err)
}
