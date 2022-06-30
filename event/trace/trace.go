package trace

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-saas/kit/event"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type SpanKind string

const (
	KindProducer SpanKind = "PRODUCER"
	KindConsumer SpanKind = "CONSUMER"
)

var (
	defaultTracer = otel.Tracer("event/middleware")
	propagator    = propagation.NewCompositeTextMapPropagator(tracing.Metadata{}, propagation.Baggage{}, propagation.TraceContext{})
)

type GetTracer interface {
	GetTracer() trace.Tracer
}
type Reporter interface {
	ReportSpanAttr() []attribute.KeyValue
}

func getTracerOrDefault(ctx context.Context) trace.Tracer {
	if r, ok := event.FromConsumerContext(ctx); ok {
		if getter, ok := r.(GetTracer); ok {
			return getter.GetTracer()
		}
	}
	if r, ok := event.FromProducerContext(ctx); ok {
		if getter, ok := r.(GetTracer); ok {
			return getter.GetTracer()
		}
	}

	return defaultTracer
}

func getAttr(ctx context.Context) []attribute.KeyValue {
	if r, ok := event.FromConsumerContext(ctx); ok {
		if getter, ok := r.(Reporter); ok {
			return getter.ReportSpanAttr()
		}
	}
	if r, ok := event.FromProducerContext(ctx); ok {
		if getter, ok := r.(Reporter); ok {
			return getter.ReportSpanAttr()
		}
	}
	return nil
}

func Send() event.ProducerMiddlewareFunc {
	return func(next event.HandlerOf[any]) event.HandlerOf[any] {
		return event.HandlerFuncOf[any](func(ctx context.Context, e interface{}) (err error) {
			var events []event.Event
			if ee, ok := e.(event.Event); ok {
				events = append(events, ee)
			} else if es, ok := e.([]event.Event); ok {
				events = append(events, es...)
			}

			ctx, span := getTracerOrDefault(ctx).Start(ctx, string(KindProducer))
			defer func() {
				if err != nil {
					span.RecordError(err)
					span.SetStatus(codes.Error, err.Error())
				} else {
					span.SetStatus(codes.Ok, "OK")
				}
				span.End()
			}()
			//span should set attr
			spanContext := span.SpanContext()
			//https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/semantic_conventions/messaging.md#apache-kafka
			attWithTopic := append(
				getAttr(ctx),
				attribute.String("span.otel.kind", string(KindProducer)),
				attribute.String("messaging.message_id", spanContext.SpanID().String()),
			)
			//extract header from ctx
			var header = propagation.HeaderCarrier{}
			propagator.Inject(ctx, header)

			//set header for each event
			for _, ee := range events {
				attWithTopic = append(attWithTopic, attribute.String("messaging.kafka.message_key", ee.Key()))
				for _, k := range header.Keys() {
					ee.Header().Set(k, header.Get(k))
				}
			}
			span.SetAttributes(attWithTopic...)
			err = next.Process(ctx, e)
			return
		})

	}
}

// Receive receive middleware
func Receive() event.ConsumerMiddlewareFunc {
	return func(next event.ConsumerHandler) event.ConsumerHandler {
		return event.ConsumerHandlerFunc(func(ctx context.Context, e event.Event) (err error) {
			//recover ctx
			header := propagation.HeaderCarrier{}
			for _, k := range e.Header().Keys() {
				header.Set(k, e.Header().Get(k))
			}
			ctx = propagator.Extract(ctx, header)
			msgKey := e.Key()

			ctx, span := getTracerOrDefault(ctx).Start(ctx, fmt.Sprintf("%s", msgKey))

			//https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/semantic_conventions/messaging.md#apache-kafka
			attWithTopic := append(
				getAttr(ctx),
				attribute.String("span.otel.kind", string(KindConsumer)),
			)
			if len(msgKey) > 0 {
				attWithTopic = append(attWithTopic, attribute.String("messaging.kafka.message_key", msgKey))
			}

			defer func() {
				if err != nil {
					span.RecordError(err)
					span.SetStatus(codes.Error, err.Error())
				} else {
					span.SetStatus(codes.Ok, "OK")
				}
				span.End()
			}()

			span.SetAttributes(attWithTopic...)

			err = next.Process(ctx, e)
			return

		})
	}
}
