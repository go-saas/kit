package server

import (
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
)

// Recovery wrap kratos recovery with handler
func Recovery() middleware.Middleware {
	return recovery.Recovery()
}
