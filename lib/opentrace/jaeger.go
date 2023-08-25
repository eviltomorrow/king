package opentrace

import (
	"context"
	"time"

	"github.com/eviltomorrow/king/lib/buildinfo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

var (
	OtelDSN = "127.0.0.1:4317"
	tracer  trace.Tracer
)

func init() {
	tracer = otel.Tracer(buildinfo.AppName)
}

func WithTracer() trace.Tracer {
	return tracer
}

func InitTraceProvider() (func() error, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(buildinfo.AppName),
			semconv.ServiceVersionKey.String(buildinfo.MainVersion),
			semconv.DeploymentEnvironmentKey.String("RELEASE"),
		),
	)
	if err != nil {
		return nil, err
	}

	exporter, err := otlptrace.New(ctx, otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(OtelDSN),
		otlptracegrpc.WithDialOption(grpc.WithBlock()),
	))
	if err != nil {
		return nil, err
	}

	processor := sdktrace.NewBatchSpanProcessor(exporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(processor),
	)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(tracerProvider)

	return func() error {
		if err := processor.Shutdown(context.Background()); err != nil {
			return err
		}

		if err := exporter.Shutdown(context.Background()); err != nil {
			return err
		}
		return tracerProvider.Shutdown(context.Background())
	}, nil
}
