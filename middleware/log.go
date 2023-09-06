package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/trace"
)

func TracingLogMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		start := time.Now()
		err := next(c)
		traceID := trace.SpanFromContext(ctx).SpanContext().TraceID().String()
		stop := time.Now()

		var errorMsg string
		if err != nil {
			errorMsg = err.Error()
		}

		c.Logger().Infoj(map[string]interface{}{
			"time":       stop.Format(time.RFC3339),
			"trace_id":   traceID,
			"id":         c.Response().Header().Get(echo.HeaderXRequestID),
			"remote_ip":  c.RealIP(),
			"host":       c.Request().Host,
			"method":     c.Request().Method,
			"uri":        c.Request().RequestURI,
			"user_agent": c.Request().UserAgent(),
			"status":     c.Response().Status,
			"error":      errorMsg,
			"latency":    stop.Sub(start).Seconds(),
			"bytes_in":   c.Request().ContentLength,
			"bytes_out":  c.Response().Size,
		})

		return err
	}
}
