package handlers

import (
	"io"
	"net/http"
	"strconv"

	"github.com/bubaew95/yandex-go-learn/internal/interfaces"
	"github.com/go-chi/chi/v5"
)

const randomStringLength = 8

type ShortenerHandler struct {
	service interfaces.ShortenerServiceInterface
}

func NewShortenerHandler(s interfaces.ShortenerServiceInterface) *ShortenerHandler {
	return &ShortenerHandler{
		service: s,
	}
}

func (s ShortenerHandler) CreateURL(res http.ResponseWriter, req *http.Request) {
	responseData, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	body := string(responseData)
	if body == "" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	genID := s.service.GenerateID(body, randomStringLength)
	url := s.service.GenerateResponseURL(genID)

	res.WriteHeader(http.StatusCreated)
	res.Header().Set("content-type", "text/plain")
	res.Header().Set("content-length", strconv.Itoa(len(url)))

	res.Write([]byte(url))
}

func (s *ShortenerHandler) GetURL(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")

	url, ok := s.service.GetURLByID(id)
	if !ok {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	res.Header().Set("Location", url)
	res.WriteHeader(http.StatusTemporaryRedirect)
}
