package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
)

var urlStore = make(map[string]string)

func shortenUrl(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()
		shortURL := fmt.Sprintf("%x", md5.Sum([]byte(body)))[:8]
		urlStore[shortURL] = string(body)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("http://localhost:8080/" + shortURL))
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func redirect(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		shortURL := r.PathValue("id")
		if urlStore[shortURL] != "" {
			w.Header().Set("Location", urlStore[shortURL])
			w.WriteHeader(http.StatusTemporaryRedirect)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, shortenUrl)
	mux.HandleFunc(`/{id}`, redirect)

	err := http.ListenAndServe(`:8080`, mux)

	if err != nil {
		panic(err)
	}
}
