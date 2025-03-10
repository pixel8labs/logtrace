package log

import (
	"context"
	"os"

	"github.com/rs/zerolog"

	"github.com/pixel8labs/logtrace/trace"
)

type Fields map[string]any

type Logger struct {
	logger        zerolog.Logger
	serviceName   string
	fieldsToScrub map[string]struct{}
}

var (
	// Initialize default logger, to support older integration.
	logger = Logger{
		logger:        zerolog.New(os.Stdout).With().Timestamp().Logger(),
		serviceName:   os.Getenv("SERVICE_NAME"),
		fieldsToScrub: map[string]struct{}{},
	}
)

func Debug(ctx context.Context, context Fields, message string, args ...any) {
	appendDefaultFields(
		ctx,
		logger.logger.Debug().Interface("context", logger.ScrubFields(context)),
	).Msgf(message, args...)
}

func Info(ctx context.Context, context Fields, message string, args ...any) {
	appendDefaultFields(
		ctx,
		logger.logger.Info().Interface("context", logger.ScrubFields(context)),
	).Msgf(message, args...)
}

func Warn(ctx context.Context, context Fields, message string, args ...any) {
	appendDefaultFields(
		ctx,
		logger.logger.Warn().Interface("context", logger.ScrubFields(context)),
	).Msgf(message, args...)
}

func Error(ctx context.Context, err error, context Fields, message string, args ...any) {
	appendDefaultFields(
		ctx,
		logger.logger.Error().Stack().Err(err).Interface("context", logger.ScrubFields(context)),
	).Msgf(message, args...)
}

func Fatal(ctx context.Context, err error, context Fields, message string, args ...any) {
	appendDefaultFields(
		ctx,
		logger.logger.Fatal().Stack().Err(err).Interface("context", logger.ScrubFields(context)),
	).Msgf(message, args...)
}

func Panic(ctx context.Context, err error, context Fields, message string, args ...any) {
	appendDefaultFields(
		ctx,
		logger.logger.Panic().Stack().Err(err).Interface("context", logger.ScrubFields(context)),
	).Msgf(message, args...)
}

func appendDefaultFields(ctx context.Context, event *zerolog.Event) *zerolog.Event {
	event = event.Str("service", logger.serviceName)
	event = appendTraceId(ctx, event)

	return event
}

func appendTraceId(ctx context.Context, event *zerolog.Event) *zerolog.Event {
	traceId, spanId := trace.TraceIdAndSpanIdFromContext(ctx)
	if traceId != "" {
		event.Str("trace_id", traceId)
	}
	if spanId != "" {
		event.Str("span_id", spanId)
	}

	return event
}
