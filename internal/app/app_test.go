package app

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
	"github.com/Hordevcom/URLShortener/internal/files"
	"github.com/Hordevcom/URLShortener/internal/middleware/logging"
	"github.com/Hordevcom/URLShortener/internal/storage"
	"github.com/go-chi/chi/v5"
)

func TestRedirect(t *testing.T) {
	logger := logging.NewLogger()
	conf := config.NewConfig()
	file := files.NewFile(conf, logger)
	m1 := storage.NewStorage(*file)
	m1.Set("abc123", "https://example.com")
	app := &App{storage: m1}

	testRequest := func(shortURL string) *http.Request {
		req := httptest.NewRequest("GET", "/"+shortURL, nil)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", shortURL)

		ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
		return req.WithContext(ctx)
	}

	t.Run("valid short URL redirects", func(t *testing.T) {
		req := testRequest("abc123")
		rr := httptest.NewRecorder()

		app.Redirect(rr, req)

		if rr.Code != http.StatusTemporaryRedirect {
			t.Errorf("expected status %d, got %d", http.StatusTemporaryRedirect, rr.Code)
		}

		expectedLocation := "https://example.com"
		if loc := rr.Header().Get("Location"); loc != expectedLocation {
			t.Errorf("expected Location header %s, got %s", expectedLocation, loc)
		}
	})

	t.Run("invalid short URL returns 400", func(t *testing.T) {
		req := testRequest("invalid")
		rr := httptest.NewRecorder()

		app.Redirect(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}

		expectedBody := "URL not found\n"
		if rr.Body.String() != expectedBody {
			t.Errorf("expected response body %q, got %q", expectedBody, rr.Body.String())
		}
	})
}

func TestShortenURL(t *testing.T) {
	logger := logging.NewLogger()
	conf := config.NewConfig()
	file := files.NewFile(conf, logger)
	m1 := storage.NewStorage(*file)
	m1.Set("abc123", "https://example.com")
	app := &App{
		storage: m1,
		config:  conf,
	}

	t.Run("successful URL shortening", func(t *testing.T) {
		originalURL := "https://example.com"
		reqBody := []byte(originalURL)

		req := httptest.NewRequest("POST", "/shorten", bytes.NewReader(reqBody))
		rr := httptest.NewRecorder()

		app.ShortenURL(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d", http.StatusCreated, rr.Code)
		}

		expectedShort := fmt.Sprintf("%x", md5.Sum(reqBody))[:8]
		expectedResponse := fmt.Sprintf("http://localhost:8080/%s", expectedShort)

		if rr.Body.String() != expectedResponse {
			t.Errorf("expected response body %q, got %q", expectedResponse, rr.Body.String())
		}

		if storedURL, exists := m1.Get(expectedShort); !exists || storedURL != originalURL {
			t.Errorf("expected storage to contain %q, but got %q", originalURL, storedURL)
		}
	})

	t.Run("empty body returns 400", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/shorten", bytes.NewReader([]byte{}))
		rr := httptest.NewRecorder()

		app.ShortenURL(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}

		expectedResponse := "url param required\n"
		if rr.Body.String() != expectedResponse {
			t.Errorf("expected response body %q, got %q", expectedResponse, rr.Body.String())
		}
	})
}

type RequestPayload struct {
	URL string `json:"url"`
}

func TestShortenURLJSON(t *testing.T) {
	logger := logging.NewLogger()
	conf := config.NewConfig()
	file := files.NewFile(conf, logger)
	app := &App{
		storage:     storage.NewStorage(*file),
		JSONStorage: *storage.NewJSONStorage(),
		config:      config.Config{Host: "http://localhost"},
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
				expectedURL := app.config.Host + "/" + expectedShort
				if res.Result != expectedURL {
					t.Errorf("expected result %s, got %s", expectedURL, res.Result)
				}
			}
		})
	}
}
