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
	"google.golang.org/grpc"
)

var (
	OtelAgentAddr = "127.0.0.1:4317"
)

func InitTraceProvider() (func() error, error) {
	ctx := context.Background()

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
		otlptracegrpc.WithEndpoint(OtelAgentAddr),
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
		cxt, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		if err := processor.Shutdown(cxt); err != nil {
			return err
		}
		if err := exporter.Shutdown(cxt); err != nil {
			return err
		}
		return tracerProvider.Shutdown(cxt)
	}, nil
}
