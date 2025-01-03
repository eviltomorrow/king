// Copyright 2019 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package zlog

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapEncodingName is the encoder name registered in zap
var ZapEncodingName = "zlog"

func DebugCtx(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx != nil {
		var traceId = trace.SpanFromContext(ctx).SpanContext().TraceID().String()
		fields = append(fields, zap.String("traceId", traceId))
	}
	ll().Debug(msg, fields...)
}

func InfoCtx(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx != nil {
		var traceId = trace.SpanFromContext(ctx).SpanContext().TraceID().String()
		fields = append(fields, zap.String("traceId", traceId))
	}
	ll().Info(msg, fields...)
}

func WarnCtx(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx != nil {
		var traceId = trace.SpanFromContext(ctx).SpanContext().TraceID().String()
		fields = append(fields, zap.String("traceId", traceId))
	}
	ll().Warn(msg, fields...)
}

func ErrorCtx(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx != nil {
		var traceId = trace.SpanFromContext(ctx).SpanContext().TraceID().String()
		fields = append(fields, zap.String("traceId", traceId))
	}
	ll().Error(msg, fields...)
}

func PanicCtx(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx != nil {
		var traceId = trace.SpanFromContext(ctx).SpanContext().TraceID().String()
		fields = append(fields, zap.String("traceId", traceId))
	}
	ll().Panic(msg, fields...)
}

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Debug(msg string, fields ...zap.Field) {
	ll().Debug(msg, fields...)
}

// Info logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Info(msg string, fields ...zap.Field) {
	ll().Info(msg, fields...)
}

// Warn logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Warn(msg string, fields ...zap.Field) {
	ll().Warn(msg, fields...)
}

// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Error(msg string, fields ...zap.Field) {
	ll().Error(msg, fields...)
}

// Panic logs a message at PanicLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then panics, even if logging at PanicLevel is disabled.
func Panic(msg string, fields ...zap.Field) {
	ll().Panic(msg, fields...)
}

// Fatal logs a message at FatalLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is
// disabled.
func Fatal(msg string, fields ...zap.Field) {
	ll().Fatal(msg, fields...)
}

// With creates a child logger and adds structured context to it.
// Fields added to the child don't affect the parent, and vice versa.
//
// Deprecated: With should not add caller skip, since it's not a logging function.
// Please use log.L().With instead. With is kept for compatibility.
// See https://github.com/pingcap/log/issues/32 for more details.
func With(fields ...zap.Field) *zap.Logger {
	return L().WithOptions(zap.AddCallerSkip(1)).With(fields...)
}

// SetLevel alters the logging level.
func SetLevel(l zapcore.Level) {
	globalProperties.Load().(*ZapProperties).Level.SetLevel(l)
}

// GetLevel gets the logging level.
func GetLevel() zapcore.Level {
	return globalProperties.Load().(*ZapProperties).Level.Level()
}
