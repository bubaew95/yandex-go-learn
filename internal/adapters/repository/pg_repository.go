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

func (p PgRepository) Close() error {
	return p.db.Close()
}

func (p PgRepository) Ping() error {
	return p.db.Ping()
}

func (p PgRepository) SetURL(id string, url string) {
	p.db.ExecContext(context.Background(),
		"INSERT INTO shortener (id, url) VALUES($1, $2)",
		id, url)
}

func (p PgRepository) GetURLByID(id string) (string, bool) {
	var url string

	row := p.db.QueryRowContext(context.Background(),
		"SELECT url FROM shortener WHERE id = $1", id)
	err := row.Scan(&url)
	if err != nil {
		return "", false
	}

	return url, true
}

func (p PgRepository) GetAllURL() map[string]string {
	return nil
}
