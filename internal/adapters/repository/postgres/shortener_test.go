package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bubaew95/yandex-go-learn/internal/adapters/constants"
	"github.com/bubaew95/yandex-go-learn/internal/core/model"
	"github.com/bubaew95/yandex-go-learn/pkg/crypto"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSetURL(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := ShortenerRepository{db: db}

	t.Run("Success added", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), crypto.KeyUserID, "1")

		mock.ExpectExec(`INSERT INTO shortener`).
			WithArgs("124f", "https://local.site", "1").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = repo.SetURL(ctx, "124f", "https://local.site")
		require.NoError(t, err)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Unique constraint violation", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), crypto.KeyUserID, "2")

		pgErr := &pgconn.PgError{Code: pgerrcode.UniqueViolation}

		mock.ExpectExec(`INSERT INTO shortener`).
			WithArgs("124f", "https://local.site", "2").
			WillReturnError(pgErr)

		err = repo.SetURL(ctx, "124f", "https://local.site")
		require.ErrorIs(t, err, constants.ErrUniqueIndex)

		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestShortenerRepository_GetURLByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		id        string
		mockRow   *sqlmock.Rows
		mockError error
		wantURL   string
		wantErr   error
	}{
		{
			name:    "found and not deleted",
			id:      "123",
			mockRow: sqlmock.NewRows([]string{"url", "is_deleted"}).AddRow("https://site.com", false),
			wantURL: "https://site.com",
			wantErr: nil,
		},
		{
			name:    "found but deleted",
			id:      "456",
			mockRow: sqlmock.NewRows([]string{"url", "is_deleted"}).AddRow("https://site.com", true),
			wantErr: constants.ErrIsDeleted,
		},
		{
			name:      "not found",
			id:        "789",
			mockError: sql.ErrNoRows,
			wantErr:   sql.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()

			repo := ShortenerRepository{db: db}
			ctx := context.Background()

			if tt.mockRow != nil {
				mock.ExpectQuery(`SELECT url, is_deleted FROM shortener WHERE id = \$1`).
					WithArgs(tt.id).
					WillReturnRows(tt.mockRow)
			} else {
				mock.ExpectQuery(`SELECT url, is_deleted FROM shortener WHERE id = \$1`).
					WithArgs(tt.id).
					WillReturnError(tt.mockError)
			}

			got, err := repo.GetURLByID(ctx, tt.id)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantURL, got)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestShortenerRepository_GetURLByOriginalURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		url       string
		mockRow   *sqlmock.Rows
		mockError error
		wantID    string
		wantFound bool
	}{
		{
			name:      "found",
			url:       "https://site.com",
			mockRow:   sqlmock.NewRows([]string{"id", "url"}).AddRow("abc", "https://site.com"),
			wantID:    "abc",
			wantFound: true,
		},
		{
			name:      "not found",
			url:       "https://missing.com",
			mockError: sql.ErrNoRows,
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()

			repo := ShortenerRepository{db: db}
			ctx := context.Background()

			if tt.mockRow != nil {
				mock.ExpectQuery(`SELECT id, url FROM shortener WHERE url = \$1`).
					WithArgs(tt.url).
					WillReturnRows(tt.mockRow)
			} else {
				mock.ExpectQuery(`SELECT id, url FROM shortener WHERE url = \$1`).
					WithArgs(tt.url).
					WillReturnError(tt.mockError)
			}

			id, ok := repo.GetURLByOriginalURL(ctx, tt.url)
			assert.Equal(t, tt.wantID, id)
			assert.Equal(t, tt.wantFound, ok)

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestShortenerRepository_GetURLSByUserID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		userID    string
		mockRows  *sqlmock.Rows
		mockError error
		wantMap   map[string]string
	}{
		{
			name:   "multiple urls",
			userID: "1",
			mockRows: sqlmock.NewRows([]string{"id", "url"}).
				AddRow("id1", "http://1").
				AddRow("id2", "http://2"),
			wantMap: map[string]string{"id1": "http://1", "id2": "http://2"},
		},
		{
			name:      "query error",
			userID:    "2",
			mockError: errors.New("db fail"),
			wantMap:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()

			repo := ShortenerRepository{db: db}
			ctx := context.Background()

			if tt.mockRows != nil {
				mock.ExpectQuery(`SELECT id, url FROM shortener WHERE user_id = \$1`).
					WithArgs(tt.userID).
					WillReturnRows(tt.mockRows)
			} else {
				mock.ExpectQuery(`SELECT id, url FROM shortener WHERE user_id = \$1`).
					WithArgs(tt.userID).
					WillReturnError(tt.mockError)
			}

			result, err := repo.GetURLSByUserID(ctx, tt.userID)

			if tt.mockError != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantMap, result)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestShortenerRepository_InsertURLs(t *testing.T) {
	t.Parallel()

	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := ShortenerRepository{db: db}
	ctx := context.WithValue(context.Background(), crypto.KeyUserID, "user-1")

	mock.ExpectBegin()

	stmt := mock.ExpectPrepare(`INSERT INTO shortener`)

	stmt.ExpectExec().
		WithArgs("abc", "http://1", "user-1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	stmt.ExpectExec().
		WithArgs("def", "http://2", "user-1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err := repo.InsertURLs(ctx, []model.ShortenerURLMapping{
		{CorrelationID: "abc", OriginalURL: "http://1"},
		{CorrelationID: "def", OriginalURL: "http://2"},
	})
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestShortenerRepository_DeleteUserURLS(t *testing.T) {
	t.Parallel()

	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := ShortenerRepository{db: db}
	ctx := context.Background()

	mock.ExpectBegin()

	stmt := mock.ExpectPrepare(`UPDATE shortener SET is_deleted = true`)

	stmt.ExpectExec().
		WithArgs("u1", "id1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	stmt.ExpectExec().
		WithArgs("u1", "id2").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err := repo.DeleteUserURLS(ctx, []model.URLToDelete{
		{ShortLink: "id1", UserID: "u1"},
		{ShortLink: "id2", UserID: "u1"},
	})
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
