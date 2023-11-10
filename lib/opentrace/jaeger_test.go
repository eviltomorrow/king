package opentrace

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/eviltomorrow/king/lib/buildinfo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var name = "jaeger-test"

func TestDemo(t *testing.T) {
	buildinfo.AppName = name
	OtelDSN = "otel-collector:4317"
	destroy, err := InitTraceProvider()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		// time.Sleep(3 * time.Second)
		destroy()
	}()

	ctx, span := otel.Tracer(name).Start(context.Background(), "Main")

	span.SetAttributes(attribute.String("ip", "192.168.33.10"))
	defer span.End()

	f1(ctx)
	// time.Sleep(2 * time.Second)
}

func f1(ctx context.Context) {
	ctx, span := otel.Tracer(name).Start(ctx, "F1()")
	defer span.End()

	span.AddEvent("lock begin")
	time.Sleep(1 * time.Second)
	span.AddEvent("lock end")
	span.RecordError(fmt.Errorf("failure"))
	span.SetStatus(codes.Error, "Failure")
	// time.Sleep(1 * time.Second)
	span1 := trace.SpanFromContext(ctx)
	id := span1.SpanContext().TraceID().String()
	fmt.Println(id)
}
