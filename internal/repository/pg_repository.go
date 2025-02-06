package repository

import (
	"context"
	"database/sql"

	"github.com/bubaew95/yandex-go-learn/config"
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
			id VARCHAR(10) PRIMARY KEY,
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

func (pr PgRepository) Close() error {
	return pr.db.Close()
}

func (pr PgRepository) Ping() error {
	return pr.db.Ping()
}

func (s PgRepository) SetURL(id string, url string) {
	s.db.ExecContext(context.Background(),
		"INSERT INTO shortener (id, url) VALUES($1, $2)",
		id, url)
}

func (s PgRepository) GetURLByID(id string) (string, bool) {
	var url string

	row := s.db.QueryRowContext(context.Background(),
		"SELECT url FROM shortener WHERE id = $1", id)
	err := row.Scan(&url)
	if err != nil {
		return "", false
	}

	return url, true
}

func (s PgRepository) GetAllURL() map[string]string {
	return nil
}
