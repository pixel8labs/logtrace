package log

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type Fields map[string]any

type Logger struct {
	logger      zerolog.Logger
	serviceName string
	env         string
}

var (
	// Initialize default logger, to support older integration.
	logger = Logger{
		logger:      zerolog.New(os.Stdout).With().Timestamp().Logger(),
		serviceName: os.Getenv("SERVICE_NAME"),
		env:         os.Getenv("APP_ENV"),
	}
)

func Debug(ctx context.Context, context Fields, message string, args ...any) {
	appendDefaultFields(ctx, logger.logger.Debug().Interface("context", context)).Msgf(message, args...)
}

func Info(ctx context.Context, context Fields, message string, args ...any) {
	appendDefaultFields(ctx, logger.logger.Info().Interface("context", context)).Msgf(message, args...)
}

func Warn(ctx context.Context, context Fields, message string, args ...any) {
	appendDefaultFields(ctx, logger.logger.Warn().Interface("context", context)).Msgf(message, args...)
}

func Error(ctx context.Context, err error, context Fields, message string, args ...any) {
	appendDefaultFields(ctx, logger.logger.Error().Stack().Err(err).Interface("context", context)).Msgf(message, args...)
}

func Fatal(ctx context.Context, err error, context Fields, message string, args ...any) {
	appendDefaultFields(ctx, logger.logger.Fatal().Stack().Err(err).Interface("context", context)).Msgf(message, args...)
}

func Panic(ctx context.Context, err error, context Fields, message string, args ...any) {
	appendDefaultFields(ctx, logger.logger.Panic().Stack().Err(err).Interface("context", context)).Msgf(message, args...)
}

func appendDefaultFields(ctx context.Context, event *zerolog.Event) *zerolog.Event {
	event = event.Str("service", logger.serviceName)
	event = event.Str("env", logger.env)
	event = appendTraceId(ctx, event)

	return event
}

func appendTraceId(ctx context.Context, event *zerolog.Event) *zerolog.Event {
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
