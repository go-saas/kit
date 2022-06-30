package kafka

import (
	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"strings"
)

// TraceSend send middleware

var (
	tracer = otel.Tracer("event/kafka")

	// These are based on the spec, which was reachable as of 2020-05-15
	// https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/semantic_conventions/messaging.md
	fixedAttrsFun = func(brokers []string) []attribute.KeyValue {
		return []attribute.KeyValue{
			attribute.String("messaging.destination_kind", "topic"),
			attribute.String("messaging.system", "kafka"),
			attribute.String("net.transport", "IP.TCP"),
			attribute.String("messaging.url", strings.Join(brokers, ","))}
	}
)

func (s *Producer) GetTracer() trace.Tracer {
	return tracer
}

func (s *Producer) ReportSpanAttr() []attribute.KeyValue {
	fixedAttrs := fixedAttrsFun(s.address)
	fixedAttrs = append(fixedAttrs,
		attribute.String("messaging.destination", s.topic),
	)
	return fixedAttrs
}

func (k *Consumer) GetTracer() trace.Tracer {
	return tracer
}
func (k *Consumer) ReportSpanAttr() []attribute.KeyValue {
	fixedAttrs := fixedAttrsFun(k.address)
	fixedAttrs = append(fixedAttrs,
		attribute.String("messaging.kafka.consumer_group", k.group),
	)
	return fixedAttrs
}
