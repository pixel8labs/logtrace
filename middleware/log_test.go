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
			name:          "Large body based on Content-Length",
			contentType:   "application/json",
			content:       strings.Repeat("A", MaxBodySize+1),
			contentLength: MaxBodySize + 1,
			expectSkip:    true,
			expectedMsg:   "Skipping body logging: Request body too large",
		},
		{
			name:          "Large body with missing Content-Length",
			contentType:   "application/json",
			content:       strings.Repeat("A", MaxBodySize+1),
			contentLength: -1,
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
