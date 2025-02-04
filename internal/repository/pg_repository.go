package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/bubaew95/yandex-go-learn/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PgRepository struct {
	db *sql.DB
}

func NewPgRepository(ctg *config.Config) *PgRepository {
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		ctg.DataBaseDSN,
		os.Getenv("PG_USER"),
		os.Getenv("PG_PASSWORD"),
		os.Getenv("PG_DB"),
	)

	return &PgRepository{
		db: dbConnect(connStr),
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pr.db.PingContext(ctx); err != nil {
		return err
	}

	return nil
}
