package log

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
)

type initConfig struct {
	writer io.Writer
	// externalWriter is to write to external resource, e.g. DataDog.
	externalWriter io.Writer
	// fieldsToScrub is a list of fields that should be scrubbed from the logs.
	fieldsToScrub []string
}

type InitOptFn func(config *initConfig)

// WithFieldsToScrub sets the fields that should be scrubbed from the logs.
// The fields are case-insensitive.
func WithFieldsToScrub(fields []string) InitOptFn {
	return func(config *initConfig) {
		config.fieldsToScrub = fields
	}
}

const dir = "/shared-logs"
const path = "app.log"

func Init(serviceName string, env string, opts ...InitOptFn) {
	cfg := &initConfig{
		writer:         os.Stdout,
		externalWriter: nil,
	}
	for _, opt := range opts {
		opt(cfg)
	}

	// Convert fields to scrub to map for faster lookup.
	fieldsToScrub := map[string]struct{}{}
	for _, field := range cfg.fieldsToScrub {
		// Use lowercase to make it case-insensitive.
		fieldsToScrub[strings.ToLower(field)] = struct{}{}
	}
	file, err := os.OpenFile(filepath.Join(dir, path), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	l := zerolog.New(file).With().Timestamp().Logger()

	logger = Logger{
		logger:        l,
		serviceName:   serviceName,
		fieldsToScrub: fieldsToScrub,
	}
}
