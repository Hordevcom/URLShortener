package app

import (
	"github.com/Hordevcom/URLShortener/internal/config"
	"github.com/Hordevcom/URLShortener/internal/files"
	"github.com/Hordevcom/URLShortener/internal/storage/pg"

	"github.com/Hordevcom/URLShortener/internal/storage"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type App struct {
	Storage     storage.Storage
	Config      config.Config
	JSONStorage storage.JSONStorage
	file        files.File
	Pg          *pg.PGDB
}

func NewApp(storage storage.Storage, config config.Config, JSONStorage storage.JSONStorage, file files.File, pg *pg.PGDB) *App {
	app := &App{Storage: storage, Config: config, file: file, Pg: pg}
	return app
}
