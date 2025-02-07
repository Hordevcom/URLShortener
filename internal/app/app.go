package app

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/Hordevcom/URLShortener/internal/config"
	"github.com/Hordevcom/URLShortener/internal/storage"
	"github.com/go-chi/chi/v5"
)

type App struct {
	storage storage.Storage
	config  config.Config
}

func NewApp(storage storage.Storage, config config.Config) *App {
	return &App{storage: storage, config: config}
}

func (a *App) ShortenURL(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}

	if len(body) == 0 {
		http.Error(w, "url param required", http.StatusBadRequest)
		return
	}
	_, err = url.ParseRequestURI(string(body))

	if err != nil {
		http.Error(w, "Correct url required", http.StatusBadRequest)
	}
	shortURL := fmt.Sprintf("%x", md5.Sum(body))[:8]
	a.storage.Set(shortURL, string(body))

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s/%s", a.config.Host, shortURL)
}

func (a *App) Redirect(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "id")
	if originalURL, exists := a.storage.Get(shortURL); exists {
		http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
	} else {
		http.Error(w, "URL not found", http.StatusBadRequest)
	}
}
