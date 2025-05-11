package handlers

import "net/http"

// DBPing осуществляет ping до базы данных
func (h *ShortenHandler) DBPing(w http.ResponseWriter, r *http.Request) {
	err := h.DB.Ping(r.Context())

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
