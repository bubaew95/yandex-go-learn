package service

import (
	"github.com/bubaew95/yandex-go-learn/internal/interfaces"
)

type PgService struct {
	repo interfaces.PgRepositoryInterface
}

func NewPgService(repo interfaces.PgRepositoryInterface) *PgService {
	return &PgService{
		repo: repo,
	}
}

func (ps PgService) Ping() error {
	return ps.repo.Ping()
}
