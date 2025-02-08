package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/bubaew95/yandex-go-learn/internal/adapters/logger"
	"github.com/bubaew95/yandex-go-learn/internal/core/model"
	"github.com/bubaew95/yandex-go-learn/internal/core/ports"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

const randomStringLength = 8

var (
	CannotDecodeJSON = "cannot decode request JSON body"
	CannotEncodeJSON = "error encoding response"
	URLNotFound      = "url not found by id"
	ErrorInsertBatch = "error insert urls by batch"
)

type ShortenerHandler struct {
	service ports.ShortenerServiceInterface
}

func NewShortenerHandler(s ports.ShortenerServiceInterface) *ShortenerHandler {
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

	url := s.service.GenerateURL(req.Context(), body, randomStringLength)

	res.WriteHeader(http.StatusCreated)
	res.Header().Set("content-type", "text/plain")
	res.Header().Set("content-length", strconv.Itoa(len(url)))

	res.Write([]byte(url))
}

func (s *ShortenerHandler) GetURL(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")

	url, ok := s.service.GetURLByID(req.Context(), id)
	if !ok {
		logger.Log.Debug(URLNotFound)
		res.WriteHeader(http.StatusNotFound)
		return
	}

	logger.Log.Info(fmt.Sprintf("Get Url by ID: %s", url))

	res.Header().Set("Location", url)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func (s *ShortenerHandler) AddNewURL(res http.ResponseWriter, req *http.Request) {
	var requestBody model.ShortenerRequest
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&requestBody); err != nil {
		logger.Log.Debug(CannotDecodeJSON, zap.Error(err))
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	url := s.service.GenerateURL(req.Context(), requestBody.URL, randomStringLength)

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)

	responseModel := model.ShortenerResponse{
		Result: url,
	}

	enc := json.NewEncoder(res)
	if err := enc.Encode(responseModel); err != nil {
		logger.Log.Debug(CannotEncodeJSON, zap.Error(err))
		return
	}
}

func (s ShortenerHandler) Ping(w http.ResponseWriter, r *http.Request) {
	if err := s.service.Ping(); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s ShortenerHandler) Batch(w http.ResponseWriter, r *http.Request) {
	var batchURLMapping []model.ShortenerURLMapping

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&batchURLMapping); err != nil {
		logger.Log.Debug(CannotDecodeJSON, zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	items, err := s.service.InsertURLs(r.Context(), batchURLMapping)
	if err != nil {
		logger.Log.Debug(ErrorInsertBatch, zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(items); err != nil {
		logger.Log.Debug(CannotEncodeJSON, zap.Error(err))
		return
	}
}
