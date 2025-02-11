package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/bubaew95/yandex-go-learn/internal/interfaces"
	"github.com/bubaew95/yandex-go-learn/internal/logger"
	"github.com/bubaew95/yandex-go-learn/internal/models"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
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
		logger.Log.Debug("body is empty")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	url := s.service.GenerateURL(body, randomStringLength)

	res.WriteHeader(http.StatusCreated)
	res.Header().Set("content-type", "text/plain")
	res.Header().Set("content-length", strconv.Itoa(len(url)))

	res.Write([]byte(url))
}

func (s *ShortenerHandler) GetURL(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")

	url, ok := s.service.GetURLByID(id)
	if !ok {
		logger.Log.Debug("url not found by id")
		res.WriteHeader(http.StatusNotFound)
		return
	}

	logger.Log.Info(fmt.Sprintf("Get Url by ID: %s", url))

	res.Header().Set("Location", url)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func (s *ShortenerHandler) AddNewURL(res http.ResponseWriter, req *http.Request) {
	var requestBody models.ShortenerRequest
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&requestBody); err != nil {
		logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	url := s.service.GenerateURL(requestBody.URL, randomStringLength)

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)

	responseModel := models.ShortenerResponse{
		Result: url,
	}

	enc := json.NewEncoder(res)
	if err := enc.Encode(responseModel); err != nil {
		logger.Log.Debug("error encoding response", zap.Error(err))
		return
	}
}
