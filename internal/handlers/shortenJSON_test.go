package handlers

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Hordevcom/URLShortener/internal/config"
	"github.com/Hordevcom/URLShortener/internal/storage"
)

type RequestPayload struct {
	URL string `json:"url"`
}

func TestShortenURLJSON(t *testing.T) {
	m1 := storage.NewMapStorage()
	m1.Set(context.Background(), "abc123", "https://example.com", 0)
	app := &ShortenHandler{
		Storage:     storage.NewMapStorage(),
		JSONStorage: *storage.NewJSONStorage(),
		Config:      config.Config{Host: "http://localhost"},
	}

	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedResult string
	}{
		{
			name:           "Valid URL",
			requestBody:    `{"url":"https://example.com"}`,
			expectedStatus: http.StatusCreated,
			expectedResult: "https://example.com",
		},
		{
			name:           "Invalid JSON",
			requestBody:    `{"url":}`,
			expectedStatus: http.StatusBadRequest,
			expectedResult: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			app.ShortenURLJSON(w, req)
			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			if resp.StatusCode == http.StatusCreated {
				var res Response
				if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
					t.Errorf("failed to decode response: %v", err)
				}

				expectedShort := fmt.Sprintf("%x", md5.Sum([]byte(tt.expectedResult)))[:8]
				expectedURL := app.Config.Host + "/" + expectedShort
				if res.Result != expectedURL {
					t.Errorf("expected result %s, got %s", expectedURL, res.Result)
				}
			}
		})
	}
}
