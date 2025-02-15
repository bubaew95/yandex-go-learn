package handlers

import (
	"fmt"
	"net/http"

	"github.com/bubaew95/yandex-go-learn/internal/adapters/logger"
	"github.com/bubaew95/yandex-go-learn/internal/core/ports"
)

type UserHandler struct {
	service ports.UserService
}

func NewUserHandler(s ports.UserService) *UserHandler {
	return &UserHandler{
		service: s,
	}
}

func (s UserHandler) GetUserURLS(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("user_id")
	if err != nil {
		logger.Log.Debug("user_id cookie not found")
	}

	fmt.Println(cookie)
}
