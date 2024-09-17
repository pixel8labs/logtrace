package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pixel8labs/logtrace/log"
	restmiddleware "github.com/pixel8labs/logtrace/middleware"
	"github.com/pixel8labs/logtrace/trace"
)

func main() {
	// Init logger.
	log.Init("logtrace-example", "local", log.WithPrettyPrint())

	// Init tracer.
	trace.InitTracer()

	e := echo.New()

	e.Use(
		// Trace middleware comes first so that the logger has the trace_id and span_id.
		restmiddleware.TracingMiddleware("logtrace-example"),
		restmiddleware.Logger(),
	)

	e.GET("/healthcheck", func(c echo.Context) error {
		log.Info(c.Request().Context(), nil, "Healthcheck received")
		return c.String(http.StatusOK, "OK")
	})

	if err := e.Start("localhost:8080"); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(context.Background(), err, nil, "Failed to start server")
		}
	}
}
