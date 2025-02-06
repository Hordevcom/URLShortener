package routes

import (
	"github.com/Hordevcom/URLShortener/internal/app"
	"github.com/go-chi/chi/v5"
)

func RegisterRouters(app app.App) *chi.Mux {
	router := chi.NewRouter()
	router.Post("/", app.ShortenURL)
	router.Get("/{id}", app.Redirect)

	return router
}
