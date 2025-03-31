package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Hordevcom/URLShortener/internal/storage"
	"github.com/go-chi/chi/v5"
)

func TestRedirect(t *testing.T) {
	m1 := storage.NewMapStorage()
	m1.Set(context.Background(), "abc123", "https://example.com", 0)
	app := &ShortenHandler{Storage: m1}

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

		if rr.Code != http.StatusGone {
			t.Errorf("expected status %d, got %d", http.StatusGone, rr.Code)
		}

		expectedBody := "URL not found\n"
		if rr.Body.String() != expectedBody {
			t.Errorf("expected response body %q, got %q", expectedBody, rr.Body.String())
		}
	})
}
