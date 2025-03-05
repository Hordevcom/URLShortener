package routes

import (
	"github.com/Hordevcom/URLShortener/internal/app"
	"github.com/Hordevcom/URLShortener/internal/middleware/compress"
	"github.com/Hordevcom/URLShortener/internal/middleware/logging"
	"github.com/go-chi/chi/v5"
)

func NewRouter(app app.App) *chi.Mux {
	router := chi.NewRouter()

	router.Use(logging.WithLogging)
	router.With(compress.DecompressMiddleware).Post("/", app.ShortenURL)
	router.With(compress.DecompressMiddleware).Post("/api/shorten", app.ShortenURLJSON)
	router.With(compress.DecompressMiddleware).Post("/api/shorten/batch", app.BatchShortenURL)
	router.Get("/ping", app.DBPing)
	router.With(compress.CompressMiddleware).Get("/{id}", app.Redirect)

	return router
}
