package log

import (
	"io"
	"net/http"
	"os"
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

func WithPrettyPrint() InitOptFn {
	return func(config *initConfig) {
		config.writer = zerolog.ConsoleWriter{Out: os.Stdout}
	}
}

func WithDataDogWriter(ddApiKey string, ddBaseUrl string, httpClient *http.Client) InitOptFn {
	return func(config *initConfig) {
		config.externalWriter = NewDataDogWriter(ddApiKey, ddBaseUrl, httpClient)
	}
}

// WithFieldsToScrub sets the fields that should be scrubbed from the logs.
// The fields are case-insensitive.
func WithFieldsToScrub(fields []string) InitOptFn {
	return func(config *initConfig) {
		config.fieldsToScrub = fields
	}
}

func Init(serviceName string, env string, opts ...InitOptFn) {
	cfg := &initConfig{
		writer:         os.Stdout,
		externalWriter: nil,
	}
	for _, opt := range opts {
		opt(cfg)
	}

	writer := cfg.writer
	if cfg.externalWriter != nil {
		writer = io.MultiWriter(writer, cfg.externalWriter)
	}

	// Convert fields to scrub to map for faster lookup.
	fieldsToScrub := map[string]struct{}{}
	for _, field := range cfg.fieldsToScrub {
		// Use lowercase to make it case-insensitive.
		fieldsToScrub[strings.ToLower(field)] = struct{}{}
	}

	l := zerolog.New(writer).With().Timestamp().Logger()

	logger = Logger{
		logger:        l,
		serviceName:   serviceName,
		env:           env,
		fieldsToScrub: fieldsToScrub,
	}
}
