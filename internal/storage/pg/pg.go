package pg

import (
	"context"
	"errors"

	"github.com/Hordevcom/URLShortener/internal/config"
	// "github.com/Hordevcom/URLShortener/internal/storage"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type PGDB struct {
	config config.Config
	logger zap.SugaredLogger
	// storage storage.Storage
}

func NewPGDB(config config.Config, logger zap.SugaredLogger) *PGDB {
	return &PGDB{config: config, logger: logger} //storage: storage
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
