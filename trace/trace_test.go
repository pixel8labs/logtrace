package trace_test

import (
	"context"
	"testing"

	"github.com/pixel8labs/logtrace/trace"
	"github.com/stretchr/testify/assert"
)

func init() {
	trace.InitTracer()
}

func TestInjectTraceToMapAndExtractTraceFromMap(t *testing.T) {
	// Given a context with a trace.
	ctx, span := trace.StartSpan(context.Background(), "test", "test-inject-extract")
	defer span.End()

	mapStringToString := make(map[string]string)
	originalTraceId, originalSpanId := trace.TraceIdAndSpanIdFromContext(ctx)
	assert.NotEmpty(t, originalTraceId)
	assert.NotEmpty(t, originalSpanId)

	// When we inject the trace to a map.
	trace.InjectTraceToMap(ctx, mapStringToString)

	// And extract the trace from the map.
	newCtx := trace.ExtractTraceFromMap(context.Background(), mapStringToString)

	// Then the trace id and span id should be the same.
	newTraceId, newSpanId := trace.TraceIdAndSpanIdFromContext(newCtx)
	assert.Equal(t, originalTraceId, newTraceId)
	assert.Equal(t, originalSpanId, newSpanId)
}
