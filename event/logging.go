package event

import (
	"context"
	klog "github.com/go-kratos/kratos/v2/log"
)

// Logging logging errors
func Logging(logger klog.Logger) ConsumerMiddlewareFunc {
	return func(next ConsumerHandler) ConsumerHandler {
		return ConsumerHandlerFunc(func(ctx context.Context, event Event) error {
			err := next.Process(ctx, event)
			if err != nil {
				_ = klog.WithContext(ctx, logger).Log(klog.LevelError,
					klog.DefaultMessageKey, err.Error(),
					"event", event.Key())
			} else {
				_ = klog.WithContext(ctx, logger).Log(klog.LevelInfo,
					"event", event.Key())
			}
			return err
		})
	}
}
