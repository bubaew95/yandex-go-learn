package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/bubaew95/yandex-go-learn/config"
	"github.com/bubaew95/yandex-go-learn/internal/core/model"
	"github.com/bubaew95/yandex-go-learn/internal/core/ports"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type ShortenerRepository struct {
	db *sql.DB
}

func NewShortenerRepository(ctg config.Config) (*ShortenerRepository, error) {
	db := dbConnect(ctg.DataBaseDSN)
	err := createTable(db)

	if err != nil {
		return nil, err
	}

	return &ShortenerRepository{
		db: db,
	}, nil
}

func createTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS shortener (
			id VARCHAR(100) PRIMARY KEY,
			url VARCHAR(1024)
		);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_original_url ON shortener (url);
	`)

	return err
}

func (p ShortenerRepository) Close() error {
	return p.db.Close()
}

func (p ShortenerRepository) Ping() error {
	return p.db.Ping()
}

func (p ShortenerRepository) SetURL(ctx context.Context, id string, url string) error {
	_, err := p.db.ExecContext(ctx,
		"INSERT INTO shortener (id, url) VALUES($1, $2)",
		id, url)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			err = ports.ErrUniqueIndex
		}
	}

	return err
}

func (p ShortenerRepository) GetURLByID(ctx context.Context, id string) (string, bool) {
	var url string

	row := p.db.QueryRowContext(ctx,
		"SELECT url FROM shortener WHERE id = $1", id)
	err := row.Scan(&url)
	if err != nil {
		return "", false
	}

	return url, true
}

func (p ShortenerRepository) GetURLByOriginalURL(ctx context.Context, originalURL string) (string, bool) {
	var (
		id  string
		url string
	)

	row := p.db.QueryRowContext(ctx,
		"SELECT id, url FROM shortener WHERE url = $1", originalURL)
	err := row.Scan(&id, &url)
	if err != nil {
		return "", false
	}

	return id, true
}

func (p ShortenerRepository) GetAllURL(ctx context.Context) map[string]string {
	return nil
}

func (p ShortenerRepository) InsertURLs(ctx context.Context, urls []model.ShortenerURLMapping) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	fmt.Println(urls)
	defer tx.Rollback()

	smtp, err := tx.PrepareContext(ctx, "INSERT INTO shortener (id, url) VALUES($1, $2) ON CONFLICT DO NOTHING")
	if err != nil {
		fmt.Println("test")
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
