package service

import (
	"context"

	"github.com/bubaew95/yandex-go-learn/internal/core/ports"
)

type UserService struct {
	repo ports.UserRepositoryInterface
}

func NewUserService(r ports.UserRepositoryInterface) *UserService {
	return &UserService{
		repo: r,
	}
}

func (u UserService) GetUserURLS(ctx context.Context, id string) {

}
