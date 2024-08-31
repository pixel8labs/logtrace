package log

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

type DataDogWriter struct {
	// apiKey is the DataDog api key.
	apiKey string
	// baseUrl is the DataDog base url.
	baseUrl    string
	httpClient *http.Client
}

// NewDataDogWriter creates a new DataDogWriter.
// apiKey is the DataDog api key.
// baseUrl is the DataDog base url, e.g. https://http-intake.logs.datadoghq.com/v1/input.
// If httpClient is nil, a default http client will be used.
func NewDataDogWriter(apiKey string, baseUrl string, httpClient *http.Client) *DataDogWriter {
	if httpClient == nil {
		httpClient = defaultHttpClient()
	}
	return &DataDogWriter{
		apiKey:     apiKey,
		baseUrl:    baseUrl,
		httpClient: httpClient,
	}
}

// Write writes the log message to DataDog.
func (w *DataDogWriter) Write(p []byte) (int, error) {
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/api/v2/logs?dd-api-key=%s", w.baseUrl, w.apiKey),
		bytes.NewBuffer(p),
	)
	if err != nil {
		return 0, fmt.Errorf("DataDogWriter http.NewRequest: %w", err)
	}

	httpRes, err := w.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("DataDogWriter httpClient.Do: %w", err)
	}
	if httpRes.StatusCode != http.StatusAccepted {
		return 0, fmt.Errorf("DataDogWriter unexpected status code from DataDog: %d", httpRes.StatusCode)
	}
	
	return len(p), nil
}

func defaultHttpClient() *http.Client {
	return &http.Client{
		Timeout: 1 * time.Second,
	}
}
