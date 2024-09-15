package restmiddleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pixel8labs/logtrace/log"
)

type customContext struct {
	echo.Context
}

// JSON override echo's c.JSON
func (c customContext) JSON(code int, payload interface{}) error {
	request := c.Request()
	err := c.Context.JSON(code, payload)

	// Log outgoing response
	resCtx := map[string]interface{}{
		"status":  c.Response().Status,
		"headers": c.Response().Header(),
		"body":    payload,
	}

	log.Info(c.Request().Context(), log.Fields{
		"request":  c.Get("request"),
		"response": resCtx,
	}, "Outgoing response: %s %s",
		request.Method,
		request.RequestURI,
	)

	return err
}

func getObject(rawData []byte) (interface{}, bool) {
	var object interface{}
	err := json.Unmarshal(rawData, &object)
	if err != nil {
		return nil, false
	}
	return object, true
}

func getRequestBody(req *http.Request) (interface{}, string, bool) {
	var reqBody []byte
	if req.Body != nil {
		reqBody, _ = io.ReadAll(req.Body)
	}
	req.Body = io.NopCloser(bytes.NewBuffer(reqBody)) // Reset

	object, ok := getObject(reqBody)
	return object, string(reqBody), ok
}

func Logger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			request := c.Request()

			// Log incoming request.
			reqCtx := map[string]interface{}{
				"method":  request.Method,
				"url":     request.RequestURI,
				"query":   c.QueryParams(),
				"headers": request.Header,
			}
			if body, rawString, ok := getRequestBody(request); ok {
				reqCtx["body"] = body
			} else {
				reqCtx["body"] = rawString
			}

			log.Info(request.Context(), log.Fields{
				"request": reqCtx,
			}, "Incoming request: %s %s",
				request.Method,
				request.RequestURI,
			)

			c.Set("request", reqCtx)

			return next(customContext{c})
		}
	}
}
