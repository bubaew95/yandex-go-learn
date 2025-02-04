package repository

import (
	"database/sql"

	"github.com/bubaew95/yandex-go-learn/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PgRepository struct {
	db *sql.DB
}

func NewPgRepository(ctg *config.Config) *PgRepository {
	return &PgRepository{
		db: dbConnect(ctg.DataBaseDSN),
	}
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
