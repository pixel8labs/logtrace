package log

import (
	"io"
	"net/http"
	"os"

	"github.com/rs/zerolog"
)

type initConfig struct {
	writer io.Writer
	// externalWriter is to write to external resource, e.g. DataDog.
	externalWriter io.Writer
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

	l := zerolog.New(writer).With().Timestamp().Logger()

	logger = Logger{
		logger:      l,
		serviceName: serviceName,
		env:         env,
	}
}
