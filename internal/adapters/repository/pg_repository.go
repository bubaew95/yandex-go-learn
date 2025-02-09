package repository

import (
	"context"
	"database/sql"

	"github.com/bubaew95/yandex-go-learn/config"
	"github.com/bubaew95/yandex-go-learn/internal/core/model"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PgRepository struct {
	db *sql.DB
}

func NewPgRepository(ctg config.Config) (*PgRepository, error) {
	db := dbConnect(ctg.DataBaseDSN)
	err := createTable(db)

	if err != nil {
		return nil, err
	}

	return &PgRepository{
		db: db,
	}, nil
}

func createTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS shortener (
			id VARCHAR(100) PRIMARY KEY,
			url VARCHAR(1024)
		)
	`)

	return err
}

func dbConnect(connStr string) *sql.DB {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		panic(err)
	}

	return db
}

func (p PgRepository) Close() error {
	return p.db.Close()
}

func (p PgRepository) Ping() error {
	return p.db.Ping()
}

func (p PgRepository) SetURL(ctx context.Context, id string, url string) error {
	_, err := p.db.ExecContext(ctx,
		"INSERT INTO shortener (id, url) VALUES($1, $2)",
		id, url)

	return err
}

func (p PgRepository) GetURLByID(ctx context.Context, id string) (string, bool) {
	var url string

	row := p.db.QueryRowContext(ctx,
		"SELECT url FROM shortener WHERE id = $1", id)
	err := row.Scan(&url)
	if err != nil {
		return "", false
	}

	return url, true
}

func (p PgRepository) GetAllURL(ctx context.Context) map[string]string {
	return nil
}

func (p PgRepository) InsertURLs(ctx context.Context, urls []model.ShortenerURLMapping) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	smtp, err := tx.PrepareContext(ctx, "INSERT INTO shortener (id, url) VALUES($1, $2) ON CONFLICT DO NOTHING")
	if err != nil {
		return err
	}
	defer smtp.Close()

	for _, v := range urls {
		_, err := smtp.ExecContext(ctx, v.CorrelationID, v.OriginalURL)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
