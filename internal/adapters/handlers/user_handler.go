package handlers

import (
	"net/http"

	"github.com/bubaew95/yandex-go-learn/internal/core/ports"
)

type UserHandler struct {
	service ports.UserServiceInterface
}

func NewUserHandler(s ports.UserServiceInterface) *UserHandler {
	return &UserHandler{
		service: s,
	}
}

func (s UserHandler) GetUserURLS(w http.ResponseWriter, r *http.Request) {

}
