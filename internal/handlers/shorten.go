package handlers

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"

	"github.com/Hordevcom/URLShortener/internal/config"
	"github.com/Hordevcom/URLShortener/internal/middleware/jwtgen"
	"github.com/Hordevcom/URLShortener/internal/storage"
	"github.com/Hordevcom/URLShortener/internal/storage/pg"
)

// ShortenHandler структура для хранения переменных
type ShortenHandler struct {
	Storage     storage.Storage
	Config      config.Config
	JSONStorage storage.JSONStorage
	DB          pg.PGDB
	DeleteCh    chan string
}

// NewShortenHandler конструктор для ShortenHandler
func NewShortenHandler(storage storage.Storage, config config.Config, JSONStorage storage.JSONStorage, db pg.PGDB, deleteCh chan string) *ShortenHandler {
	return &ShortenHandler{
		Storage:     storage,
		Config:      config,
		JSONStorage: JSONStorage,
		DB:          db,
		DeleteCh:    deleteCh,
	}
}

// ShortenURL осуществляет сокращение урла, переданного с помощью text/plain
func (h *ShortenHandler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")

	if err != nil {
		fmt.Print("No token value!")
	}

	UserID := jwtgen.GetUserID(cookie.Value)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}

	if len(body) == 0 {
		http.Error(w, "url param required", http.StatusBadRequest)
		return
	}

	shortURL := fmt.Sprintf("%x", md5.Sum(body))[:8]

	ok := h.Storage.Set(r.Context(), shortURL, string(body), UserID)

	if !ok {
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintf(w, "%s/%s", h.Config.Host, shortURL)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s/%s", h.Config.Host, shortURL)
}
