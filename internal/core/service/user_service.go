package service

import (
	"context"

	"github.com/bubaew95/yandex-go-learn/internal/core/ports"
)

type UserService struct {
	repo ports.UserRepository
}

func NewUserService(r ports.UserRepository) *UserService {
	return &UserService{
		repo: r,
	}
}

func (u UserService) GetUserURLS(ctx context.Context, id string) {

}
