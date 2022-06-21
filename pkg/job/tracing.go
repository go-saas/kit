package job

import (
	"context"
	"fmt"
	"github.com/hibiken/asynq"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"net/http"
)

type SpanKind string

const (
	KindProducer SpanKind = "PRODUCER"
	KindConsumer SpanKind = "CONSUMER"
)

var (
	fixedAttrs = []attribute.KeyValue{
		attribute.String("job.system", "asynq"),
	}
	tracer     = otel.Tracer("asynq/tasks")
	propagator = propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{})
)

func TracingServer() asynq.MiddlewareFunc {
	kind := KindConsumer
	return func(h asynq.Handler) asynq.Handler {
		return asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) (err error) {
			//TODO recover context from header?
			//ctx = propagator.Extract(ctx, propagation.HeaderCarrier(t.Header))
			ctx, span := tracer.Start(ctx, fmt.Sprintf("job-%s", t.Type()))
			id, _ := asynq.GetTaskID(ctx)
			queue, _ := asynq.GetQueueName(ctx)
			maxRetry, _ := asynq.GetMaxRetry(ctx)
			retryCount, _ := asynq.GetRetryCount(ctx)
			attrs := append(
				fixedAttrs,
				attribute.String("span.otel.kind", string(kind)),
				attribute.String("job.job_id", id),
				attribute.String("job.queue", queue),
				attribute.Int("job.max_retry", maxRetry),
				attribute.Int("job.retry_count", retryCount),
			)
			span.SetAttributes(attrs...)
			defer func() {
				if err != nil {
					span.RecordError(err)
					span.SetStatus(codes.Error, err.Error())
				} else {
					span.SetStatus(codes.Ok, "OK")
				}
				span.End()
			}()
			err = h.ProcessTask(ctx, t)
			return
		})
	}
}

func SetTracingOption(ctx context.Context) asynq.Option {
	//recover header
	var header = propagation.HeaderCarrier(http.Header{})
	ctx = propagator.Extract(ctx, header)
	//TODO ability to set header
	panic("unimplemented")
}

func EnqueueWithTracing(ctx context.Context, client *asynq.Client, task *asynq.Task, opts ...asynq.Option) (taskinfo *asynq.TaskInfo, err error) {
	ctx, span := tracer.Start(ctx, fmt.Sprintf("job-%s", task.Type()))
	defer func() {
		attrs := append(
			fixedAttrs,
			attribute.String("span.otel.kind", string(KindProducer)),
		)
		if taskinfo != nil {
			attrs = append(attrs, attribute.String("job.job_id", taskinfo.ID),
				attribute.String("job.queue", taskinfo.Queue),
				attribute.Int("job.max_retry", taskinfo.MaxRetry))
		}
		span.SetAttributes(attrs...)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		} else {
			span.SetStatus(codes.Ok, "OK")
		}
		span.End()
	}()

	opts = append(opts, SetTracingOption(ctx))
	taskinfo, err = client.EnqueueContext(ctx, task, opts...)
	return
}
