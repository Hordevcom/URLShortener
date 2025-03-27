package routes

import (
	"github.com/Hordevcom/URLShortener/internal/handlers"
	"github.com/Hordevcom/URLShortener/internal/middleware/compress"
	"github.com/Hordevcom/URLShortener/internal/middleware/jwtgen"
	"github.com/Hordevcom/URLShortener/internal/middleware/logging"
	"github.com/go-chi/chi/v5"
)

func NewRouter(handler handlers.ShortenHandler) *chi.Mux {
	router := chi.NewRouter()

	router.Use(logging.WithLogging)
	router.With(compress.DecompressMiddleware, jwtgen.AuthMiddleware).Post("/", handler.ShortenURL)
	router.With(compress.DecompressMiddleware).Post("/api/shorten", handler.ShortenURLJSON)
	router.With(compress.DecompressMiddleware).Post("/api/shorten/batch", handler.BatchShortenURL)
	router.Get("/ping", handler.DBPing)
	router.Get("/api/user/urls", handler.GetUserUrls)
	router.Delete("/api/user/urls", handler.DeleteUrls)
	router.With(compress.CompressMiddleware).Get("/{id}", handler.Redirect)

	return router
}
