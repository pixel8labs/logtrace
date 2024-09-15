package restmiddleware

import (
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
)

// TracingMiddleware is a middleware that creates a new span for each incoming request.
// TODO: propagate the trace from the request headers.
func TracingMiddleware(serviceName string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tracer := otel.Tracer(serviceName)
			req := c.Request()
			ctx, span := tracer.Start(req.Context(), c.Path())

			defer span.End()

			c.SetRequest(req.WithContext(ctx))

			return next(c)
		}
	}
}
