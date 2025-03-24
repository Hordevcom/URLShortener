package handlers

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Hordevcom/URLShortener/internal/config"
	"github.com/Hordevcom/URLShortener/internal/storage"
)

func TestShortenURL(t *testing.T) {
	m1 := storage.NewMapStorage()
	m1.Set(context.Background(), "abc123", "https://example.com", 0)
	conf := config.NewConfig()

	app := &ShortenHandler{
		Storage: m1,
		Config:  conf,
	}

	t.Run("successful URL shortening", func(t *testing.T) {
		originalURL := "https://example.com"
		reqBody := []byte(originalURL)

		req := httptest.NewRequest("POST", "/shorten", bytes.NewReader(reqBody))

		req.AddCookie(&http.Cookie{
			Name:  "token",
			Value: "2",
		})
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

		if storedURL, exists := m1.Get(req.Context(), expectedShort); !exists || storedURL != originalURL {
			t.Errorf("expected storage to contain %q, but got %q", originalURL, storedURL)
		}
	})

	t.Run("empty body returns 400", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/shorten", bytes.NewReader([]byte{}))

		req.AddCookie(&http.Cookie{
			Name:  "token",
			Value: "2",
		})
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
