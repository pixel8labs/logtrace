package restmiddleware

import (
	"github.com/labstack/echo/v4"
	"github.com/pixel8labs/logtrace/trace"
)

// Tracer is a middleware that creates a new span for each incoming request.
// TODO: propagate the trace from the request headers.
func Tracer(serviceName string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			ctx, span := trace.StartSpan(req.Context(), serviceName, c.Path())

			defer span.End()

			c.SetRequest(req.WithContext(ctx))

			return next(c)
		}
	}
}
