package log

import (
	"context"
	"io"
	"os"
	"github.com/rs/zerolog"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type Fields map[string]any

var (
	logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
)

func SetOutput(w io.Writer) {
	logger = logger.Output(w)
}

func Debug(ctx context.Context, context Fields, message string, args ...any) {
	appendTraceID(ctx, logger.Debug().Interface("context", context)).Msgf(message, args...)
}

func Info(ctx context.Context, context Fields, message string, args ...any) {
	appendTraceID(ctx, logger.Info().Interface("context", context)).Msgf(message, args...)
}

func Warn(ctx context.Context, context Fields, message string, args ...any) {
	appendTraceID(ctx, logger.Warn().Interface("context", context)).Msgf(message, args...)
}

func Error(ctx context.Context, err error, context Fields, message string, args ...any) {
	appendTraceID(ctx, logger.Error().Stack().Err(err).Interface("context", context)).Msgf(message, args...)
}

func Fatal(ctx context.Context, err error, context Fields, message string, args ...any) {
	appendTraceID(ctx, logger.Fatal().Stack().Err(err).Interface("context", context)).Msgf(message, args...)
}

func Panic(ctx context.Context, err error, context Fields, message string, args ...any) {
	appendTraceID(ctx, logger.Panic().Stack().Err(err).Interface("context", context)).Msgf(message, args...)
}

func appendTraceID(ctx context.Context, event *zerolog.Event) *zerolog.Event {
	span := oteltrace.SpanFromContext(ctx)
	if span != nil {
		if span.SpanContext().HasTraceID() {
			event.Str("trace_id", span.SpanContext().TraceID().String())
		}
		if span.SpanContext().HasSpanID() {
			event.Str("span_id", span.SpanContext().SpanID().String())
		}
	}

	return event
}
