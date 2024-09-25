// Package trace is a wrapper around the OpenTelemetry tracing library.
// This provides simplified function to do something needed by Pixel8Labs (e.g. init trace, start span, extract ids).
package trace

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func InitTracer() {
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})
}

func StartSpan(ctx context.Context, serviceName string, spanName string) (context.Context, trace.Span) {
	return otel.Tracer(serviceName).Start(ctx, spanName)
}

func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

func TraceIdAndSpanIdFromContext(ctx context.Context) (string, string) {
	span := SpanFromContext(ctx)
	traceId := ""
	spanId := ""
	if span == nil {
		return traceId, spanId
	}

	if span.SpanContext().HasTraceID() {
		traceId = span.SpanContext().TraceID().String()
	}
	if span.SpanContext().HasSpanID() {
		spanId = span.SpanContext().SpanID().String()
	}

	return traceId, spanId
}

// TODO: implement start new span with trace_id
