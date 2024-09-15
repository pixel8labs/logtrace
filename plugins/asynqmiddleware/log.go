package asynqmiddleware

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/pixel8labs/logtrace/log"
)

// Logger is an asynq middleware that will log the incoming message.
// It'll also log for failure and success in processing the message.
func Logger() asynq.MiddlewareFunc {
	return func(next asynq.Handler) asynq.Handler {
		return asynq.HandlerFunc(func(ctx context.Context, task *asynq.Task) error {
			logFields := log.Fields{
				"queue":   task.Type(), // We use task type as the queue name.
				"payload": string(task.Payload()),
			}
			log.Info(ctx, logFields, "Processing queue message...")

			if err := next.ProcessTask(ctx, task); err != nil {
				log.Error(ctx, err, logFields, "Failed to process queue message")
				return err
			}

			log.Info(ctx, logFields, "Processed queue message successfully!")

			return nil
		})
	}
}
