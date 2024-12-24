package middleware

import (
	"context"
	"fmt"

	"github.com/eviltomorrow/king/lib/opentrace"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryClientOpentraceInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	ctx = extract(ctx, otel.GetTextMapPropagator())
	span1 := trace.SpanFromContext(ctx)
	// id := span1.SpanContext().TraceID().String()
	// fmt.Println(id)

	name, attr := spanInfo(info.FullMethod, peerFromCtx(ctx))

	startOpts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindServer),
		trace.WithAttributes(attr...),
	}

	ctx, span := opentrace.DefaultTracer().Start(
		trace.ContextWithRemoteSpanContext(ctx, trace.SpanContextFromContext(ctx)),
		name,
		startOpts...,
	)
	defer span.End()

	resp, err = handler(ctx, req)
	if err != nil {
		s, _ := status.FromError(err)
		statusCode, msg := serverStatus(s)
		span.SetStatus(statusCode, msg)
		span.SetAttributes(statusCodeAttr(s.Code()))
	} else {
		span.SetAttributes(statusCodeAttr(codes.OK))
	}
	return resp, err
}

// UnaryServerRecoveryInterceptor recover
func UnaryServerOpentraceInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	ctx = extract(ctx, otel.GetTextMapPropagator())
	span1 := trace.SpanFromContext(ctx)
	id := span1.SpanContext().TraceID().String()
	fmt.Println(id)

	name, attr := spanInfo(info.FullMethod, peerFromCtx(ctx))

	startOpts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindServer),
		trace.WithAttributes(attr...),
	}

	ctx, span := opentrace.DefaultTracer().Start(
		trace.ContextWithRemoteSpanContext(ctx, trace.SpanContextFromContext(ctx)),
		name,
		startOpts...,
	)
	defer span.End()

	resp, err = handler(ctx, req)
	if err != nil {
		s, _ := status.FromError(err)
		statusCode, msg := serverStatus(s)
		span.SetStatus(statusCode, msg)
		span.SetAttributes(statusCodeAttr(s.Code()))
	} else {
		span.SetAttributes(statusCodeAttr(codes.OK))
	}
	return resp, err
}

// StreamServerRecoveryInterceptor recover
func StreamServerOpentraceInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	ctx := extract(stream.Context(), otel.GetTextMapPropagator())
	name, attr := spanInfo(info.FullMethod, peerFromCtx(ctx))

	startOpts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindServer),
		trace.WithAttributes(attr...),
	}

	ctx, span := opentrace.DefaultTracer().Start(
		trace.ContextWithRemoteSpanContext(ctx, trace.SpanContextFromContext(ctx)),
		name,
		startOpts...,
	)
	defer span.End()

	err = handler(srv, &grpc_middleware.WrappedServerStream{
		ServerStream:   stream,
		WrappedContext: ctx,
	})
	if err != nil {
		s, _ := status.FromError(err)
		statusCode, msg := serverStatus(s)
		span.SetStatus(statusCode, msg)
		span.SetAttributes(statusCodeAttr(s.Code()))
	} else {
		span.SetAttributes(statusCodeAttr(codes.OK))
	}
	return err
}
