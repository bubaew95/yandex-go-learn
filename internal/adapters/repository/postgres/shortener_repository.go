package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/bubaew95/yandex-go-learn/config"
	"github.com/bubaew95/yandex-go-learn/internal/adapters/logger"
	"github.com/bubaew95/yandex-go-learn/internal/core/model"
	"github.com/bubaew95/yandex-go-learn/internal/core/ports"
	"github.com/bubaew95/yandex-go-learn/pkg/crypto"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
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
			url VARCHAR(1024),
			user_id VARCHAR(255),
			is_deleted BOOLEAN DEFAULT FALSE
		);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_original_url ON shortener (url);
		CREATE INDEX IF NOT EXISTS idx_user_id ON shortener (user_id);
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
	userID := ctx.Value(crypto.KeyUserID)

	logger.Log.Debug("SetURL", zap.Any("user_id", userID))
	_, err := p.db.ExecContext(ctx,
		"INSERT INTO shortener (id, url, user_id) VALUES($1, $2, $3)",
		id, url, userID)

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
		"SELECT url FROM shortener WHERE id = $1 and is_deleted = false", id)
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
	defer tx.Rollback()

	userID := ctx.Value(crypto.KeyUserID)

	smtp, err := tx.PrepareContext(ctx, "INSERT INTO shortener (id, url, user_id) VALUES($1, $2, $3) ON CONFLICT DO NOTHING")
	if err != nil {
		fmt.Println("test")
		return err
	}
	defer smtp.Close()

	for _, v := range urls {
		_, err := smtp.ExecContext(ctx, v.CorrelationID, v.OriginalURL, userID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
func (p ShortenerRepository) GetURLSByUserID(ctx context.Context, userID string) (map[string]string, error) {
	rows, err := p.db.QueryContext(ctx, "SELECT id, url FROM shortener WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make(map[string]string)
	for rows.Next() {
		var (
			ID  string
			URL string
		)
		err := rows.Scan(&ID, &URL)
		if err != nil {
			return nil, err
		}

		items[ID] = URL
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (p ShortenerRepository) DeleteUserURLS(ctx context.Context, items []model.URLToDelete) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, "UPDATE shortener SET is_deleted = true WHERE is_deleted = false AND user_id = $1 AND id = $2")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, item := range items {
		_, err := stmt.ExecContext(ctx, item.UserID, item.ShortLink)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
