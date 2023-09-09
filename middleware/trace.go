package restmiddleware

import (
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
)

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