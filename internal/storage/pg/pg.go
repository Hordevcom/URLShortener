package pg

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"

	"github.com/Hordevcom/URLShortener/internal/config"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

type PGDB struct {
	config config.Config
	logger zap.SugaredLogger
}

func NewPGDB(config config.Config, logger zap.SugaredLogger) *PGDB {

	return &PGDB{config: config, logger: logger}
}

func (p *PGDB) ConnectToDB() (*pgxpool.Pool, error) {
	db, err := pgxpool.New(context.Background(), p.config.DatabaseDsn) //sql.Open("pgx", p.config.DatabaseDsn)

	if err != nil {
		p.logger.Errorw("Problem with connecting to db: ", err)
		return nil, err
	}

	err = db.Ping(context.Background())

	if err != nil {
		p.logger.Errorw("Problem with ping to db: ", err)
		return nil, err
	}

	p.logger.Infow("Connecting and ping to db successful")
	return db, nil
}

// func (p *PGDB) CreateTable(db *sql.DB) {
// 	createTableSQL := `
// 	CREATE TABLE IF NOT EXISTS urls (
// 		short_url TEXT NOT NULL PRIMARY KEY,
// 		original_url TEXT NOT NULL
// 	);`

// 	_, err := db.Exec(createTableSQL)

// 	if err != nil {
// 		p.logger.Errorw("Cannot create table: ", err)
// 		return
// 	}
// }

func (p *PGDB) Get(shortURL string) (string, bool) {
	var origURL string
	db, err := p.ConnectToDB()

	if err != nil {
		p.logger.Errorw("Error to connect to db: ", err)
		return "", false
	}
	defer db.Close()

	row := db.QueryRow(context.Background(), `SELECT original_url FROM urls WHERE short_url = $1`, shortURL)
	row.Scan(&origURL)

	if origURL == "" {
		return "", false
	}
	return origURL, true
}

func (p *PGDB) Set(shortURL, originalURL string) bool {
	query := `INSERT INTO urls (short_url, original_url)
	 VALUES ($1, $2) ON CONFLICT (short_url) DO NOTHING`

	db, err := p.ConnectToDB()

	if err != nil {
		p.logger.Errorw("Error to connect to db: ", err)
		return false
	}

	result, err := db.Exec(context.Background(), query, shortURL, originalURL)

	if rows := result.RowsAffected(); rows == 0 {
		return false
	}
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return false
		}
	}

	return true
}

func InitMigrations(conf config.Config, logger zap.SugaredLogger) {
	logger.Infow("Start migrations")
	db, err := sql.Open("pgx", conf.DatabaseDsn)

	if err != nil {
		logger.Fatalw("Error with connection to DB: ", err)
	}

	defer db.Close()

	wd, _ := os.Getwd()
	migrationsPath := filepath.Join(filepath.Dir(filepath.Dir(wd)), "internal", "storage", "migrations")

	err = goose.Up(db, migrationsPath)
	if err != nil {
		logger.Fatalw("Error with migrations: ", err)
	}
}
