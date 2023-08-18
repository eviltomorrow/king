package opentrace

import (
	"context"
	"testing"
	"time"

	"github.com/eviltomorrow/king/lib/buildinfo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var (
	name   = "test"
	tracer trace.Tracer
)

func TestDemo(t *testing.T) {
	buildinfo.AppName = name
	OtelAgentAddr = "192.168.233.10:4317"
	destroy, err := InitTraceProvider()
	if err != nil {
		t.Fatal(err)
	}
	defer destroy()

	tracer = otel.Tracer(name)
	ctx, span := tracer.Start(context.Background(), "TestMain")
	defer span.End()

	f1(ctx)
	time.Sleep(3 * time.Second)
}

func f1(ctx context.Context) {
	ctx, span := otel.Tracer(name).Start(ctx, "F1")

	defer span.End()

	f2(ctx)
	time.Sleep(1 * time.Second)
}

func f2(ctx context.Context) {
	_, span := otel.Tracer(name).Start(ctx, "f2")
	defer span.End()

	time.Sleep(2 * time.Second)

}
