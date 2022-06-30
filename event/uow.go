package event

import (
	"context"
	"github.com/go-saas/uow"
)

//ConsumerUow wrap handler into a unit of work (transaction)
func ConsumerUow(uowMgr uow.Manager) ConsumerMiddlewareFunc {
	return func(handler ConsumerHandler) ConsumerHandler {
		return ConsumerHandlerFunc(func(ctx context.Context, event Event) error {
			return uowMgr.WithNew(ctx, func(ctx context.Context) error {
				return handler.Process(ctx, event)
			})
		})
	}
}
