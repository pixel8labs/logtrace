package asynqmiddleware

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/pixel8labs/logtrace/trace"
)

// Tracer is an asynq middleware that will start a new span for each request.
// TODO: propagate the trace from the incoming message.
func Tracer(serviceName string) asynq.MiddlewareFunc {
	return func(next asynq.Handler) asynq.Handler {
		return asynq.HandlerFunc(func(ctx context.Context, task *asynq.Task) error {
			ctx, span := trace.StartSpan(ctx, serviceName, task.Type())
			defer span.End()
			return next.ProcessTask(ctx, task)
		})
	}
}
