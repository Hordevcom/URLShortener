package app

import (
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Hordevcom/URLShortener/internal/config"
	"github.com/Hordevcom/URLShortener/internal/files"
	"github.com/Hordevcom/URLShortener/internal/storage"
	"github.com/go-chi/chi/v5"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type App struct {
	storage     storage.Storage
	config      config.Config
	JSONStorage storage.JSONStorage
	file        files.File
}

type Response struct {
	Result string `json:"result"`
}

func NewApp(storage storage.Storage, config config.Config, JSONStorage storage.JSONStorage, file files.File) *App {
	return &App{storage: storage, config: config, file: file}
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

	shortURL := fmt.Sprintf("%x", md5.Sum(body))[:8]

	if _, exist := a.storage.Get(shortURL); !exist {
		a.storage.Set(shortURL, string(body))

		a.file.UpdateFile(files.JSONStruct{
			ShortURL:    shortURL,
			OriginalURL: string(body),
		})
		a.addDataToDB(shortURL, string(body))
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

	if _, exist := a.storage.Get(shortURL); !exist {
		a.storage.Set(shortURL, a.JSONStorage.Get())
		a.file.UpdateFile(files.JSONStruct{
			ShortURL:    shortURL,
			OriginalURL: a.JSONStorage.Get(),
		})
		a.addDataToDB(shortURL, a.JSONStorage.Get())

	}

	response := Response{
		Result: a.config.Host + "/" + shortURL,
	}

	JSONResponse, _ := json.Marshal(response)

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

func (a *App) ConnectToDB(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("pgx", a.config.DatabaseDsn)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed connect to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	err = db.Ping()

	if err != nil {
		http.Error(w, "Failed ping to database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *App) addDataToDB(shortURL, originalURL string) {
	db, err := sql.Open("pgx", a.config.DatabaseDsn)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	err = db.Ping()

	if err != nil {
		return
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS urls (
		short_url TEXT NOT NULL PRIMARY KEY,
		original_url TEXT NOT NULL
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		panic(err)
	}

	query := `INSERT INTO urls VALUES ($1, $2)`
	_, err = db.Exec(query, shortURL, originalURL)
	if err != nil {
		panic(err)
	}
}
