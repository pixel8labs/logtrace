package log_test

import (
	"context"

	"github.com/pixel8labs/logtrace/log"
)

func ExampleInfo_WithoutInit() {
	log.Info(context.Background(), log.Fields{"key": "value"}, "Hello, World!")
	// Example output: {"level":"info","context":{"key":"value"},"service":"","env":"","time":"2025-02-04T20:57:21+07:00","message":"Hello, World!"}
	// Can't put the actual output here because the time is dynamic.
}

func ExampleInit_WithFieldsToScrub() {
	log.Init("service-name", "development", log.WithFieldsToScrub([]string{"password"}))
	log.Info(context.Background(), log.Fields{
		"password": "shouldbescrubbed",
		"username": "name",
		"mapOfFields": map[string]any{
			"key":      "value",
			"password": "shouldbescrubbed",
		},
		"list": []any{
			struct {
				Password string
				Uname    string
			}{
				Password: "shouldbescrubbed",
				Uname:    "namenotbescrubbed",
			},
			map[int]any{
				1: map[string]any{
					"key":      true,
					"password": false,
				},
			},
		},
	}, "Hello, World!")
	// Example output: {"level":"info","context":{"list":[{"Password":"***scrubbed***","Uname":"namenotbescrubbed"},{"1":{"key":true,"password":"***scrubbed***"}}],"mapOfFields":{"key":"value","password":"***scrubbed***"},"password":"***scrubbed***","username":"name"},"service":"service-name","env":"development","time":"2025-02-04T20:55:40+07:00","message":"Hello, World!"}
	// Can't put the actual output here because the time is dynamic.
}
