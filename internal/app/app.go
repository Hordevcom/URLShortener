package app

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Hordevcom/URLShortener/internal/config"
	"github.com/Hordevcom/URLShortener/internal/files"
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

	db, err := a.pg.ConnectToDB()

	if err != nil {
		http.Error(w, "Problem with connecting to DB", http.StatusInternalServerError)
		return
	}
	var responces []ShortenResponce
	for _, req := range requests {
		shortURL := fmt.Sprintf("%x", md5.Sum([]byte(req.OriginalURL)))[:8]
		responces = append(responces, ShortenResponce{
			CorrelationID: req.CorrelationID,
			ShortURL:      shortURL,
		})

		if _, exist := a.storage.Get(shortURL); !exist {
			a.storage.Set(shortURL, req.OriginalURL)
			a.pg.AddValuesToDB(db, shortURL, req.OriginalURL)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responces)
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

func (a *App) DBPing(w http.ResponseWriter, r *http.Request) {
	db, _ := a.pg.ConnectToDB()
	defer db.Close()
	w.WriteHeader(http.StatusOK)
}

func (a *App) InitDB() {
	db, err := a.pg.ConnectToDB()

	if err != nil {
		panic(err)
	}
	a.pg.CreateTable(db)
	defer db.Close()
}

func (a *App) addDataToDB(shortURL, originalURL string) {
	db, _ := a.pg.ConnectToDB()
	defer db.Close()
	// a.pg.CreateTable(db)
	a.pg.AddValuesToDB(db, shortURL, originalURL)
}

// func (a *App) DownloadData() {
// 	db, err := a.pg.ConnectToDB()

// 	if err != nil {
// 		a.file.ReadFile(a.storage)
// 	}
// 	defer db.Close()
// 	a.pg.ReadDataFromDB(db)
// }
