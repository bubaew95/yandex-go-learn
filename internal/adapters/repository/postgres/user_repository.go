package postgres

import (
	"context"
	"database/sql"

	"github.com/bubaew95/yandex-go-learn/config"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(ctg config.Config) *UserRepository {
	db := dbConnect(ctg.DataBaseDSN)

	return &UserRepository{
		db,
	}
}

func (u UserRepository) GetUserURLS(ctx context.Context, id string) {
}
