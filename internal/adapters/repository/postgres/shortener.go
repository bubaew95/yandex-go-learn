// Package postgres реализует хранилище сокращённых URL на основе PostgreSQL.
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"

	"github.com/bubaew95/yandex-go-learn/config"
	"github.com/bubaew95/yandex-go-learn/internal/adapters/constants"
	"github.com/bubaew95/yandex-go-learn/internal/adapters/logger"
	"github.com/bubaew95/yandex-go-learn/internal/core/model"
	"github.com/bubaew95/yandex-go-learn/pkg/crypto"
)

// ShortenerRepository реализует интерфейс репозитория для работы с сокращёнными URL.
// Использует PostgreSQL как хранилище данных.
type ShortenerRepository struct {
	db *sql.DB
}

// NewShortenerRepository создаёт и инициализирует новый экземпляр ShortenerRepository,
// выполняет подключение к БД и создаёт таблицы, если их нет.
//
// Возвращает ошибку, если соединение с БД или инициализация схемы завершились неудачно.
func NewShortenerRepository(ctg config.Config) (*ShortenerRepository, error) {
	db, err := dbConnect(ctg.DataBaseDSN)
	if err != nil {
		return nil, err
	}

	err = createTable(db)

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

// Close закрывает соединение с базой данных.
func (p ShortenerRepository) Close() error {
	return p.db.Close()
}

// Ping проверяет доступность подключения к базе данных.
func (p ShortenerRepository) Ping(ctx context.Context) error {
	return p.db.PingContext(ctx)
}

// SetURL сохраняет новый сокращённый URL в базу данных.
// В случае конфликта уникального ключа возвращает ошибку ErrUniqueIndex.
func (p ShortenerRepository) SetURL(ctx context.Context, id string, url string) error {
	userID := ctx.Value(crypto.KeyUserID)

	logger.Log.Debug("SetURL", zap.Any("user_id", userID))
	_, err := p.db.ExecContext(ctx,
		"INSERT INTO shortener (id, url, user_id) VALUES($1, $2, $3)",
		id, url, userID)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			err = constants.ErrUniqueIndex
		}
	}

	return err
}

// GetURLByID возвращает оригинальный URL по его сокращённому идентификатору.
// Если запись помечена как удалённая, возвращает ошибку ErrIsDeleted.
func (p ShortenerRepository) GetURLByID(ctx context.Context, id string) (string, error) {
	var (
		url       string
		isDeleted bool
	)

	row := p.db.QueryRowContext(ctx,
		"SELECT url, is_deleted FROM shortener WHERE id = $1", id)
	err := row.Scan(&url, &isDeleted)
	if err != nil {
		return "", err
	}

	if isDeleted {
		return "", constants.ErrIsDeleted
	}

	return url, nil
}

// GetURLByOriginalURL ищет короткий ID по оригинальному URL.
// Возвращает false, если совпадение не найдено.
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

// InsertURLs добавляет список URL в БД, пропуская уже существующие записи (ON CONFLICT DO NOTHING).
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

// GetURLSByUserID возвращает все сокращённые ссылки, созданные пользователем.
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
		err = rows.Scan(&ID, &URL)
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

// DeleteUserURLS помечает указанные пользователем URL как удалённые (is_deleted = true).
func (p ShortenerRepository) DeleteUserURLS(ctx context.Context, items []model.URLToDelete) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, "UPDATE shortener SET is_deleted = true WHERE user_id = $1 and id = $2")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, item := range items {
		logger.Log.Debug("v_id", zap.String("id", item.ShortLink))
		_, err := stmt.ExecContext(ctx, item.UserID, item.ShortLink)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
