package pg

import (
	"context"
	"database/sql"

	"github.com/Hordevcom/URLShortener/internal/config"
	"github.com/Hordevcom/URLShortener/internal/storage"
	"go.uber.org/zap"
)

type PGDB struct {
	config  config.Config
	logger  zap.SugaredLogger
	storage storage.Storage
}

func NewPGDB(config config.Config, logger zap.SugaredLogger, storage storage.Storage) *PGDB {
	return &PGDB{config: config, logger: logger, storage: storage}
}

func (p *PGDB) ConnectToDB() (*sql.DB, error) {
	db, err := sql.Open("pgx", p.config.DatabaseDsn)

	if err != nil {
		p.logger.Errorw("Problem with connecting to db: ", err)
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		p.logger.Errorw("Problem with ping to db: ", err)
		return nil, err
	}

	p.logger.Infow("Connecting and ping to db successful")
	return db, nil
}

func (p *PGDB) CreateTable(db *sql.DB) {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS urls (
		short_url TEXT NOT NULL PRIMARY KEY,
		original_url TEXT NOT NULL
	);`

	_, err := db.Exec(createTableSQL)

	if err != nil {
		p.logger.Errorw("Cannot create table: ", err)
		return
	}
}

func (p *PGDB) AddValuesToDB(db *sql.DB, shortURL, originalURL string) {
	query := `INSERT INTO urls VALUES ($1, $2)`
	_, err := db.Exec(query, shortURL, originalURL)
	if err != nil {
		p.logger.Errorw("Cannot add value to table: ", err)
		return
	}
}

func (p *PGDB) ReadDataFromDB(db *sql.DB) {
	rows, err := db.QueryContext(context.Background(),
		"SELECT short_url, original_url FROM urls")

	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var shortURL string
		var origURL string

		err := rows.Scan(&shortURL, &origURL)

		if err != nil {
			p.logger.Errorw("Read row error: ", err)
		}

		p.storage.Set(shortURL, origURL)
	}
	p.logger.Infow("Read data from db completed")
}
