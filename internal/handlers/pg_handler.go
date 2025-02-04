package handlers

import (
	"fmt"
	"net/http"

	"github.com/bubaew95/yandex-go-learn/internal/interfaces"
)

type PgHandler struct {
	service interfaces.PgServiceInterface
}

func NewPgHandler(s interfaces.PgServiceInterface) *PgHandler {
	return &PgHandler{
		service: s,
	}
}

func (ps PgHandler) Ping(w http.ResponseWriter, r *http.Request) {
	if err := ps.service.Ping(); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
