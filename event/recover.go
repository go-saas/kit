package event

import (
	"context"
	"fmt"
	klog "github.com/go-kratos/kratos/v2/log"
)

type RecoverOption func(*recoverOptions)

type recoverOptions struct {
	formatter ErrFormatFunc
	logger    klog.Logger
}

type ErrFormatFunc func(ctx context.Context, err error) error

func WithErrorFormatter(f ErrFormatFunc) RecoverOption {
	return func(o *recoverOptions) {
		o.formatter = f
	}
}

func WithLogger(logger klog.Logger) RecoverOption {
	return func(o *recoverOptions) {
		o.logger = logger
	}
}

//ConsumerRecover prevent consumer from panic
func ConsumerRecover(opt ...RecoverOption) ConsumerMiddlewareFunc {
	op := recoverOptions{
		logger: klog.GetLogger(),
		formatter: func(ctx context.Context, err error) error {
			return err
		},
	}
	for _, o := range opt {
		o(&op)
	}
	logger := klog.NewHelper(op.logger)
	return func(next ConsumerHandler) ConsumerHandler {
		return ConsumerHandlerFunc(func(ctx context.Context, event Event) (err error) {
			defer func() {
				if rerr := recover(); rerr != nil {
					if rrerr, ok := rerr.(error); ok {
						wrrer := fmt.Errorf("panic recovered: %w", rrerr)
						logger.Error(wrrer)
						err = op.formatter(ctx, wrrer)
					} else {
						err = fmt.Errorf("panic recovered: %s", rerr)
						logger.Error(err)
						err = op.formatter(ctx, err)
					}
				}
			}()
			err = next.Process(ctx, event)
			if err == nil {
				return nil
			}
			return op.formatter(ctx, err)
		})
	}
}
