package main

import (
	"crypto/md5"
	"fmt"
	"net/http"
)

var urlStore = make(map[string]string)

func shortenURL(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		url := r.FormValue("url")
		shortURL := fmt.Sprintf("%x", md5.Sum([]byte(url)))[:8]
		urlStore[shortURL] = url
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("http://localhost:8080/" + shortURL))
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func redirect(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		shortUrl := r.PathValue("id")

		if urlStore[shortUrl] != "" {
			w.WriteHeader(http.StatusTemporaryRedirect)
			w.Write([]byte("Location: " + urlStore[shortUrl]))
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, shortenURL)
	mux.HandleFunc(`/{id}`, redirect)

	err := http.ListenAndServe(`:8080`, mux)

	if err != nil {
		panic(err)
	}
}
