package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/bubaew95/yandex-go-learn/internal/adapters/logger"
	"github.com/bubaew95/yandex-go-learn/internal/core/model"
	"github.com/bubaew95/yandex-go-learn/internal/core/ports"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const randomStringLength = 8

type ShortenerHandler struct {
	service ports.ShortenerService
}

func NewShortenerHandler(s ports.ShortenerService) *ShortenerHandler {
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

	url, err := s.service.GenerateURL(req.Context(), body, randomStringLength)
	if err != nil {
		if errors.Is(err, ports.ErrUniqueIndex) {
			originURL, ok := s.service.GetURLByOriginalURL(req.Context(), body)

			if ok {
				logger.Log.Debug("Duplicate", zap.String("originURL", originURL), zap.String("bodyUrl", body))

				writeByteResponse(res, http.StatusConflict, []byte(originURL))
				return
			}
		}

		logger.Log.Debug("Generate url error", zap.Error(err))
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeByteResponse(res, http.StatusCreated, []byte(url))
}

func (s *ShortenerHandler) GetURL(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")

	url, ok := s.service.GetURLByID(req.Context(), id)
	if !ok {
		logger.Log.Debug("url not found by id", zap.String("id", id))
		res.WriteHeader(http.StatusNotFound)
		return
	}

	res.Header().Set("Location", url)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func (s *ShortenerHandler) AddNewURL(res http.ResponseWriter, req *http.Request) {
	var requestBody model.ShortenerRequest

	if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
		logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	url, err := s.service.GenerateURL(req.Context(), requestBody.URL, randomStringLength)
	if err != nil {
		if errors.Is(err, ports.ErrUniqueIndex) {
			originURL, ok := s.service.GetURLByOriginalURL(req.Context(), requestBody.URL)
			if ok {
				responseModel := model.ShortenerResponse{
					Result: originURL,
				}

				writeJSONResponse(res, http.StatusConflict, responseModel, logger.Log)
				return
			}
		}

		logger.Log.Debug("Url generation error", zap.Error(err))
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseModel := model.ShortenerResponse{
		Result: url,
	}

	writeJSONResponse(res, http.StatusCreated, responseModel, logger.Log)
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
		logger.Log.Debug("cannot decode request JSON", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	items, err := s.service.InsertURLs(r.Context(), batchURLMapping)
	if err != nil {
		logger.Log.Debug("error insert urls by batch", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusCreated, items, logger.Log)
}

func (s ShortenerHandler) GetUserURLS(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("user_id")
	if err != nil {
		logger.Log.Debug("user_id cookie not found")
	}

	fmt.Println(cookie)
}

func writeJSONResponse(res http.ResponseWriter, statusCode int, data interface{}, logger *zap.Logger) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(statusCode)

	if err := json.NewEncoder(res).Encode(data); err != nil {
		logger.Debug("Cannot encode JSON", zap.Error(err))
	}
}

func writeByteResponse(res http.ResponseWriter, statusCode int, data []byte) {
	res.WriteHeader(statusCode)
	res.Header().Set("content-type", "text/plain")
	res.Header().Set("content-length", strconv.Itoa(len(data)))

	res.Write(data)
}
