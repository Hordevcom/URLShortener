package app

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Hordevcom/URLShortener/internal/config"
	"github.com/Hordevcom/URLShortener/internal/files"
	"github.com/Hordevcom/URLShortener/internal/middleware/jwtgen"
	"github.com/Hordevcom/URLShortener/internal/storage/pg"

	"github.com/Hordevcom/URLShortener/internal/storage"

	"github.com/go-chi/chi/v5"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type ShortenRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type ShortenResponce struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type ShortenOrigURLs struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type App struct {
	storage     storage.Storage
	config      config.Config
	JSONStorage storage.JSONStorage
	file        files.File
	pg          *pg.PGDB
}

type Response struct {
	Result string `json:"result"`
}

func NewApp(storage storage.Storage, config config.Config, JSONStorage storage.JSONStorage, file files.File, pg *pg.PGDB) *App {
	app := &App{storage: storage, config: config, file: file, pg: pg}
	// app.DownloadData() create bug!
	return app
}

func (a *App) BatchShortenURL(w http.ResponseWriter, r *http.Request) {
	var requests []ShortenRequest

	err := json.NewDecoder(r.Body).Decode(&requests)

	if err != nil {
		http.Error(w, "Bad JSON data", http.StatusBadRequest)
		return
	}

	if len(requests) == 0 {
		http.Error(w, "Batch cannot be empty", http.StatusBadRequest)
		return
	}

	var responces []ShortenResponce
	for _, req := range requests {
		shortURL := fmt.Sprintf("%x", md5.Sum([]byte(req.OriginalURL)))[:8]
		responces = append(responces, ShortenResponce{
			CorrelationID: req.CorrelationID,
			ShortURL:      a.config.Host + "/" + shortURL,
		})

		a.storage.Set(shortURL, req.OriginalURL, 0)

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(responces)
}

func (a *App) ShortenURL(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	UserID := 0
	token := ""

	if err != nil {
		token, _ = jwtgen.BuildJWTString()
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    token,
			HttpOnly: true,
		})

		// cookie.Value = token
		fmt.Println("Coockie set!")
		fmt.Println(token)
	}
	if cookie.Valid() == nil {
		UserID = jwtgen.GetUserID(cookie.Value)
		fmt.Println("Coockie value taken from coockie!")
	} else {
		UserID = jwtgen.GetUserID(token)
		fmt.Println("Coockie value taken from token!")
	}

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

	ok := a.storage.Set(shortURL, string(body), UserID)
	if !ok {
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintf(w, "%s/%s", a.config.Host, shortURL)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s/%s", a.config.Host, shortURL)
}

func (a *App) ShortenURLJSON(w http.ResponseWriter, r *http.Request) {
	// extract string from json
	err := json.NewDecoder(r.Body).Decode(&a.JSONStorage)

	if err != nil {
		http.Error(w, "Bad JSON url", http.StatusBadRequest)
		return
	}

	shortURL := fmt.Sprintf("%x", md5.Sum([]byte(a.JSONStorage.Get())))[:8]

	response := Response{
		Result: a.config.Host + "/" + shortURL,
	}

	JSONResponse, _ := json.Marshal(response)

	if !a.storage.Set(shortURL, a.JSONStorage.Get(), 0) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		w.Write(JSONResponse)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(JSONResponse)

}

func (a *App) Redirect(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "id")
	if originalURL, exists := a.storage.Get(shortURL); exists {
		http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
	} else {
		http.Error(w, "URL not found", http.StatusBadRequest)
	}
}

func (a *App) DBPing(w http.ResponseWriter, r *http.Request) {
	err := a.pg.Ping()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *App) GetUserUrls(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	var UserID int
	var ShorigURLs []ShortenOrigURLs

	if err != nil {
		token, _ := jwtgen.BuildJWTString()
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    token,
			HttpOnly: true,
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		fmt.Println("Token created at GetUserUrls!")
		return
	}
	fmt.Println(cookie.Value)
	if err := cookie.Valid(); err == nil {
		UserID = jwtgen.GetUserID(cookie.Value)
		fmt.Println(UserID)
		fmt.Println("UserID collected from cookie.Value")
	}

	URLs, ok := a.pg.GetWithUserID(UserID)

	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	for key, value := range URLs {
		ShorigURLs = append(ShorigURLs, ShortenOrigURLs{
			ShortURL:    a.config.Host + "/" + key,
			OriginalURL: value,
		})
	}
	var ShorigURLs1 []ShortenOrigURLs
	ShorigURLs1 = append(ShorigURLs1, ShorigURLs[len(ShorigURLs)-1])
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(ShorigURLs1)
	if err != nil {
		panic(err)
	}

}
