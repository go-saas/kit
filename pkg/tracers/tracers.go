package tracers

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-saas/kit/pkg/conf"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var DefaultPropagator = propagation.NewCompositeTextMapPropagator(tracing.Metadata{}, propagation.Baggage{}, propagation.TraceContext{})

func SetTracerProvider(ctx context.Context, cfg *conf.Tracers, name string) (func(), error) {
	if cfg != nil && cfg.Otel != nil {
		//set as otel
		conn, err := grpc.DialContext(ctx, cfg.Otel.Grpc, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return func() {

			}, err
		}
		// Set up a trace exporter
		traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
		if err != nil {
			return func() {

			}, err
		}
		// Register the trace exporter with a TracerProvider, using a batch
		// span processor to aggregate spans before export.
		bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
		tracerProvider := sdktrace.NewTracerProvider(
			// Set the sampling rate based on the parent span to 100%
			sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(1.0))),
			sdktrace.WithResource(resource.NewSchemaless(
				semconv.ServiceNameKey.String(name),
			)),
			sdktrace.WithSpanProcessor(bsp),
		)
		otel.SetTracerProvider(tracerProvider)
		return func() {
			// Shutdown will flush any remaining spans and shut down the exporter.
			err := tracerProvider.Shutdown(ctx)
			if err != nil {
				log.Error(err)
			}
		}, nil
	}
	//do not support
	return func() {

	}, nil
}
