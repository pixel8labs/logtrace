package restmiddleware

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestGetRequestBody(t *testing.T) {
	tests := []struct {
		name          string
		contentType   string
		content       string
		contentLength int64
		expectSkip    bool
		expectedMsg   string
	}{
		{
			name:          "Normal JSON request",
			contentType:   "application/json",
			content:       `{"name":"test"}`,
			contentLength: int64(len(`{"name":"test"}`)),
			expectSkip:    false,
		},
		{
			name:          "Multipart form request",
			contentType:   "multipart/form-data; boundary=12345",
			content:       "fake-multipart-data",
			contentLength: int64(len("fake-multipart-data")),
			expectSkip:    true,
			expectedMsg:   "Skipping body logging: multipart/form-data",
		},
		{
			name:          "Large request body",
			contentType:   "application/json",
			content:       strings.Repeat("A", MaxBodySize+1), // 16KB+1
			contentLength: int64(MaxBodySize + 1),
			expectSkip:    true,
			expectedMsg:   "Skipping body logging: Request body too large",
		},
		{
			name:          "Content-Length exceeds limit",
			contentType:   "application/json",
			content:       `{"test": "data"}`,
			contentLength: MaxBodySize + 1,
			expectSkip:    true,
			expectedMsg:   "Skipping body logging: Request body too large",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{
				Header: make(http.Header),
				Body:   io.NopCloser(bytes.NewBufferString(tt.content)),
			}
			req.Header.Set("Content-Type", tt.contentType)
			req.ContentLength = tt.contentLength

			body, raw, ok := getRequestBody(req)

			if tt.expectSkip {
				if ok || body != nil || raw != tt.expectedMsg {
					t.Errorf("Expected skipping, got body=%v, raw=%s", body, raw)
				}
			} else {
				if !ok || body == nil {
					t.Errorf("Expected body to be logged, but was skipped")
				}
			}
		})
	}
}
