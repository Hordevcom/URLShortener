package pg

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"runtime"

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
	db     *pgxpool.Pool
}

func NewPGDB(config config.Config, logger zap.SugaredLogger) *PGDB {
	db, err := pgxpool.New(context.Background(), config.DatabaseDsn)

	if err != nil {
		logger.Errorw("Problem with connecting to db: ", err)
		return nil
	}
	return &PGDB{config: config, logger: logger, db: db}
}

func (p *PGDB) UpdateDeleteParam(shortURLs []string) {
	query := `UPDATE urls
				SET is_deleted = TRUE
				WHERE short_url = ANY($1)`

	_, err := p.db.Exec(context.Background(), query, shortURLs)
	if err != nil {
		p.logger.Errorw("Update table error: ", err)
		return
	}
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

func (p *PGDB) Ping() error {
	err := p.db.Ping(context.Background())

	if err != nil {
		p.logger.Errorw("Problem with ping to db: ", err)
		return err
	}

	return nil
}

func (p *PGDB) Get(shortURL string) (string, bool) {
	var origURL string

	query := `SELECT original_url FROM urls WHERE short_url = $1`
	row := p.db.QueryRow(context.Background(), query, shortURL)
	row.Scan(&origURL)

	if origURL == "" {
		return "", false
	}
	return origURL, true
}

func (p *PGDB) GetWithUserID(UserID int) (map[string]string, bool) {
	var origURL string
	var shortURL string
	URLs := make(map[string]string)

	query := `SELECT original_url, short_url FROM urls WHERE user_id = $1`
	rows, err := p.db.Query(context.Background(), query, UserID)

	if err != nil {
		p.logger.Fatalw("Ошибка выполнения запроса %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&origURL, &shortURL)
		if err != nil {
			p.logger.Fatalw("Ошибка сканирования строки: %v", err)
		}

		URLs[shortURL] = origURL
	}

	if origURL == "" {
		return nil, false
	}
	return URLs, true
}

func (p *PGDB) Set(shortURL, originalURL string, userID int) bool {
	query := `INSERT INTO urls (short_url, original_url, user_id)
	 VALUES ($1, $2, $3) ON CONFLICT (short_url) DO NOTHING`

	result, err := p.db.Exec(context.Background(), query, shortURL, originalURL, userID)

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

	_, filename, _, _ := runtime.Caller(0)
	migrationsPath := filepath.Join(filepath.Dir(filename), "..", "migrations")

	err = goose.Up(db, migrationsPath)
	if err != nil {
		logger.Fatalw("Error with migrations: ", err)
	}
}
