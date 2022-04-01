package kafka

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"net/http"
	"strings"

	"github.com/Shopify/sarama"
	"go.opentelemetry.io/otel/trace"
)

type SpanKind string

const (
	KindProducer SpanKind = "PRODUCER"
	KindConsumer SpanKind = "CONSUMER"
)

type OTelInterceptor struct {
	tracer     trace.Tracer
	fixedAttrs []attribute.KeyValue
	propagator propagation.TextMapPropagator
}

// NewOTelInterceptor processes span for intercepted messages and add some
// headers with the span data.
func NewOTelInterceptor(kind SpanKind, brokers []string) *OTelInterceptor {
	oi := OTelInterceptor{}
	oi.tracer = otel.Tracer("kafka/interceptors")
	oi.propagator = propagation.NewCompositeTextMapPropagator(tracing.Metadata{}, propagation.Baggage{}, propagation.TraceContext{})
	// These are based on the spec, which was reachable as of 2020-05-15
	// https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/semantic_conventions/messaging.md
	oi.fixedAttrs = []attribute.KeyValue{
		attribute.String("messaging.destination_kind", "topic"),
		attribute.String("span.otel.kind", string(kind)),
		attribute.String("messaging.system", "kafka"),
		attribute.String("net.transport", "IP.TCP"),
		attribute.String("messaging.url", strings.Join(brokers, ",")),
	}
	return &oi
}

const (
	hasTraceHeader = "has_trace"
)

func (oi *OTelInterceptor) hasTrace(msg []*sarama.RecordHeader) bool {
	// check message hasn't been here before (retries)
	for _, h := range msg {
		if string(h.Key) == hasTraceHeader {
			return true
		}
	}
	return false
}

func (oi *OTelInterceptor) OnSend(msg *sarama.ProducerMessage) {
	h := make([]*sarama.RecordHeader, len(msg.Headers))
	for i, recordHeader := range msg.Headers {
		h[i] = &recordHeader
	}
	if oi.hasTrace(h) {
		return
	}
	ctx := context.Background()
	//recover header
	var header transport.Header = propagation.HeaderCarrier(http.Header{})
	if c, ok := msg.Metadata.(context.Context); ok {
		ctx = c
		if ts, ok := transport.FromServerContext(ctx); ok {
			header = ts.RequestHeader()
		}
		if ts, ok := transport.FromClientContext(ctx); ok {
			header = ts.RequestHeader()
		}
	}
	ctx = oi.propagator.Extract(ctx, header)
	ctx, span := oi.tracer.Start(ctx, msg.Topic)
	msg.Metadata = ctx
	defer func() { span.End() }()

	spanContext := span.SpanContext()
	//https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/semantic_conventions/messaging.md#apache-kafka
	attWithTopic := append(
		oi.fixedAttrs,
		attribute.String("messaging.destination", msg.Topic),
		attribute.String("messaging.message_id", spanContext.SpanID().String()),
	)
	if msg.Key != nil {
		if b, err := msg.Key.Encode(); err == nil {
			attWithTopic = append(attWithTopic, attribute.String("messaging.kafka.message_key", string(b)))
		}
	}
	attWithTopic = append(attWithTopic, attribute.Int("messaging.kafka.partition", int(msg.Partition)))

	span.SetAttributes(attWithTopic...)

	// remove existing partial tracing headers if exists
	noTraceHeaders := msg.Headers[:0]
	for _, h := range msg.Headers {
		key := string(h.Key)
		if key != hasTraceHeader {
			noTraceHeaders = append(noTraceHeaders, h)
		}
	}
	//header inject
	traceHeaders := []sarama.RecordHeader{
		{Key: []byte(hasTraceHeader), Value: []byte("1")},
	}

	injectHeader := http.Header{}
	oi.propagator.Inject(ctx, propagation.HeaderCarrier(injectHeader))
	for s, ss := range injectHeader {
		for _, sss := range ss {
			traceHeaders = append(traceHeaders, sarama.RecordHeader{
				Key: []byte(s), Value: []byte(sss),
			})
		}
	}

	msg.Headers = append(noTraceHeaders, traceHeaders...)
}

func (oi *OTelInterceptor) StartConsumerSpan(ctx context.Context, group string, msg *sarama.ConsumerMessage) (context.Context, trace.Span) {
	header := http.Header{}
	if oi.hasTrace(msg.Headers) {
		for _, recordHeader := range msg.Headers {
			header.Set(string(recordHeader.Key), string(recordHeader.Value))
		}
	}
	ctx = oi.propagator.Extract(ctx, propagation.HeaderCarrier(header))
	ctx, span := oi.tracer.Start(ctx, msg.Topic)

	//https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/semantic_conventions/messaging.md#apache-kafka
	attWithTopic := append(
		oi.fixedAttrs,
		attribute.String("messaging.kafka.consumer_group", group),
	)
	if msg.Key != nil {
		attWithTopic = append(attWithTopic, attribute.String("messaging.kafka.message_key", string(msg.Key)))
	}
	attWithTopic = append(attWithTopic, attribute.Int("messaging.kafka.partition", int(msg.Partition)))
	span.SetAttributes(attWithTopic...)
	return ctx, span
}
