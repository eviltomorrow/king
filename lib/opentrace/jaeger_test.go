package opentrace

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/eviltomorrow/king/lib/buildinfo"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var name = "jaeger2"

func TestDemo(t *testing.T) {
	buildinfo.AppName = name

	destroy, err := InitTraceProvider(&Config{
		DSN:            "127.0.0.1:4317",
		ConnectTimeout: 5 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		// time.Sleep(3 * time.Second)
		destroy()
	}()

	ctx, span := DefaultTracer().Start(context.Background(), "Main")

	span.SetAttributes(attribute.String("ip", "192.168.33.10"))
	defer span.End()

	f1(ctx)
	// time.Sleep(2 * time.Second)
}

func f1(ctx context.Context) {
	ctx, span := DefaultTracer().Start(ctx, "F1()")
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
