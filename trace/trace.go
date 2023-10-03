package trace

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func InitTrace(ctx context.Context) error {
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return nil
}

func Tracer(name string) trace.Tracer {
	return otel.Tracer(name)
}

func StartSpan(ctx context.Context, serviceName string, startName string) (context.Context, trace.Span) {
	return otel.GetTracerProvider().Tracer(serviceName).Start(ctx, startName)
}
