package handlers

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Result string `json:"result"`
}

func (h *ShortenHandler) ShortenURLJSON(w http.ResponseWriter, r *http.Request) {
	// extract string from json
	err := json.NewDecoder(r.Body).Decode(&h.JSONStorage)

	if err != nil {
		http.Error(w, "Bad JSON url", http.StatusBadRequest)
		return
	}

	shortURL := fmt.Sprintf("%x", md5.Sum([]byte(h.JSONStorage.Get())))[:8]

	response := Response{
		Result: h.Config.Host + "/" + shortURL,
	}

	JSONResponse, _ := json.Marshal(response)

	if !h.Storage.Set(r.Context(), shortURL, h.JSONStorage.Get(), 0) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		w.Write(JSONResponse)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(JSONResponse)

}
