package routes

import (
	"github.com/Hordevcom/URLShortener/internal/app"
	"github.com/Hordevcom/URLShortener/internal/middleware/logging"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func NewRouter(app app.App, sugar zap.SugaredLogger) *chi.Mux {
	router := chi.NewRouter()
	router.Post("/", logging.WithLogging(app.ShortenURL, sugar))
	router.Get("/{id}", logging.WithLogging(app.Redirect, sugar))

	return router
}
