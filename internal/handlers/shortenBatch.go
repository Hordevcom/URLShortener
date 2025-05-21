package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Hordevcom/URLShortener/internal/storage/pg"
)

// BatchShortenURL хандлер для обработки группы сайтов
func (h *ShortenHandler) BatchShortenURL(w http.ResponseWriter, r *http.Request) {
	var requests []pg.ShortenRequest

	err := json.NewDecoder(r.Body).Decode(&requests)

	if err != nil {
		http.Error(w, "Bad JSON data", http.StatusBadRequest)
		return
	}

	if len(requests) == 0 {
		http.Error(w, "Batch cannot be empty", http.StatusBadRequest)
		return
	}

	responces, err := h.DB.BatchShortenURL(r.Context(), requests)

	if err != nil {
		http.Error(w, "Failed to process request", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(responces)
}
