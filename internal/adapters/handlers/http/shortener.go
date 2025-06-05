package http

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/bubaew95/yandex-go-learn/config"
	"io"
	"net"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/bubaew95/yandex-go-learn/internal/adapters/constants"
	"github.com/bubaew95/yandex-go-learn/internal/adapters/logger"
	"github.com/bubaew95/yandex-go-learn/internal/core/model"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const randomStringLength = 8

// ShortenerService определяет бизнес-логику сервиса сокращения ссылок.
// Включает в себя генерацию ссылок, работу с пользователями и отложенное удаление.
//
//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name=ShortenerService --filename=servicemock_test.go --inpackage
type ShortenerService interface {
	// GenerateURL генерирует короткий URL на основе оригинального.
	GenerateURL(ctx context.Context, url string, randomStringLength int) (string, error)

	// GetURLByID возвращает оригинальный URL по его сокращённому ID.
	GetURLByID(ctx context.Context, id string) (string, error)

	// GetURLByOriginalURL возвращает ID, соответствующий оригинальному URL.
	GetURLByOriginalURL(ctx context.Context, originalURL string) (string, bool)

	// InsertURLs добавляет множество URL и возвращает их короткие представления.
	InsertURLs(ctx context.Context, urls []model.ShortenerURLMapping) ([]model.ShortenerURLResponse, error)

	// GetURLSByUserID возвращает список сокращённых URL, принадлежащих пользователю.
	GetURLSByUserID(ctx context.Context, userID string) ([]model.ShortenerURLSForUserResponse, error)

	// DeleteUserURLS помечает ссылки как удалённые.
	DeleteUserURLS(ctx context.Context, items []model.URLToDelete) error

	// ScheduleURLDeletion планирует асинхронное удаление ссылок (например, через очередь).
	ScheduleURLDeletion(ctx context.Context, items []model.URLToDelete)

	// RandStringBytes генерирует случайную строку заданной длины (обычно для ID короткой ссылки).
	RandStringBytes(n int) string

	// Ping проверяет доступность сервиса (например, для liveness-проб).
	Ping(ctx context.Context) error

	// Stats возвращает статистику
	Stats(ctx context.Context) (model.StatsRespose, error)
}

// ShortenerHandler обрабатывает HTTP-запросы, связанные с сокращением URL.
type ShortenerHandler struct {
	service ShortenerService
	config  *config.Config
}

// NewShortenerHandler возвращает новый экземпляр ShortenerHandler.
func NewShortenerHandler(s ShortenerService, cfg config.Config) *ShortenerHandler {
	return &ShortenerHandler{
		service: s,
		config:  &cfg,
	}
}

func writeJSONResponse(res http.ResponseWriter, statusCode int, data interface{}) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(statusCode)

	if err := json.NewEncoder(res).Encode(data); err != nil {
		logger.Log.Debug("Cannot encode JSON", zap.Error(err))
	}
}

func writeByteResponse(res http.ResponseWriter, statusCode int, data []byte) {
	res.WriteHeader(statusCode)
	res.Header().Set("content-type", "text/plain")
	res.Header().Set("content-length", strconv.Itoa(len(data)))

	res.Write(data)
}

// CreateURL обрабатывает HTTP POST-запрос на создание короткой ссылки.
//
// Ожидает оригинальный URL в теле запроса (как текст).
// Возвращает укороченную ссылку в случае успеха.
// Если такая ссылка уже есть — возвращает HTTP 409 и ранее созданную короткую ссылку.
func (s ShortenerHandler) CreateURL(res http.ResponseWriter, req *http.Request) {
	responseData, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	body := string(responseData)
	if body == "" {
		logger.Log.Debug("Body is empty")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	url, err := s.service.GenerateURL(req.Context(), body, randomStringLength)
	if err != nil {
		if errors.Is(err, constants.ErrUniqueIndex) {
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

// GetURL обрабатывает GET-запрос для получения оригинального URL по его короткому идентификатору.
//
// Ожидает параметр id, по которому извлекается оригинальная ссылка.
// Если ссылка найдена возврашает HTTP 307 статус и перенаправляет на оригинальную ссылку.
// Если ссылка удалена возврашает HTTP 410 статус.
// Если ссылка не найдена - возврашает HTTP 404 статус.
func (s ShortenerHandler) GetURL(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")

	url, err := s.service.GetURLByID(req.Context(), id)
	if err != nil || url == "" {
		if errors.Is(err, constants.ErrIsDeleted) {
			logger.Log.Debug("Url is deleted", zap.String("id", id))
			res.WriteHeader(http.StatusGone)
			return
		}

		logger.Log.Debug("Url not found by id", zap.String("id", id))
		res.WriteHeader(http.StatusNotFound)
		return
	}

	res.Header().Set("Location", url)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

// AddNewURL обрабатывает HTTP POST-запрос на создание короткой ссылки.
//
// Ожидает JSON данные в теое запроса.
// Возврашает HTTP 201 статус и JSON тело ответа.
// Если JSON тело запроса имеет ошибку - вовзврашается HTTP 500 ошибка.
// Если при генерации короткой ссылки возникла ошибка - возврается HTTP 500 ошибка.
// Если такая ссылка уже добавлена в базу - возврашается оригинальная ссылка из базы.
func (s ShortenerHandler) AddNewURL(res http.ResponseWriter, req *http.Request) {
	var requestBody model.ShortenerRequest

	if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
		logger.Log.Debug("Cannot decode request JSON body", zap.Error(err))
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	url, err := s.service.GenerateURL(req.Context(), requestBody.URL, randomStringLength)
	if err != nil {
		if errors.Is(err, constants.ErrUniqueIndex) {
			originURL, ok := s.service.GetURLByOriginalURL(req.Context(), requestBody.URL)
			if ok {
				responseModel := model.ShortenerResponse{
					Result: originURL,
				}

				writeJSONResponse(res, http.StatusConflict, responseModel)
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

	writeJSONResponse(res, http.StatusCreated, responseModel)
}

// Ping - обрабатывает HTTP GET-запрос на проверку подключения к БД.
//
// Если подключение успешное - возврашает HTTP 200 статус.
// Если при подключении возникла ошибка - возврашает HTTP 500 статус.
func (s ShortenerHandler) Ping(w http.ResponseWriter, r *http.Request) {
	if err := s.service.Ping(r.Context()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Batch - обрабатывает HTTP POST-запрос на создание которих ссылок.
//
// Ожидает JSON массив с которотким ID и с оригинальной ссылкой.
// Если добавление ссылок прошла успешно - возврашает HTTP 201 статус и все добавленыее ссылки.
// Если в JSON есть ошибка - возврашает HTTP 500 ошибку.
// Если при добавлении возникла ошибка - возврашает HTTP 500 статус.
func (s ShortenerHandler) Batch(w http.ResponseWriter, r *http.Request) {
	var batchURLMapping []model.ShortenerURLMapping

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&batchURLMapping); err != nil {
		logger.Log.Debug("Cannot decode request JSON", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	items, err := s.service.InsertURLs(r.Context(), batchURLMapping)
	if err != nil {
		logger.Log.Debug("Error insert urls by batch", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusCreated, items)
}

// GetUserURLS - обрабатывает HTTP GET-запрос на получение ссылок авторизованного пользователя.
//
// Если есть ссылки - возврашает HTTP 200 статус и все ссылки.
// Если ссылок нет - возврашает HTTP 204 статус.
// Если в запросе возникла ошибка возврашает HTTP 500 ошибку.
func (s ShortenerHandler) GetUserURLS(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("user_id")
	if err != nil {
		logger.Log.Debug("Cookie not found")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	items, err := s.service.GetURLSByUserID(r.Context(), cookie.Value)
	if err != nil {
		logger.Log.Debug("Get urls error", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if items == nil {
		logger.Log.Debug("User urls not found")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	writeJSONResponse(w, http.StatusOK, items)
}

// DeleteUserURLS - обрабатывает HTTP DELETE-запрос на удаление ссылок из базы.
func (s ShortenerHandler) DeleteUserURLS(w http.ResponseWriter, r *http.Request) {
	userID, err := r.Cookie("user_id")
	if err != nil || userID.Value == "" {
		logger.Log.Debug("Cookie not found")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var deleteItems []string
	if err := json.NewDecoder(r.Body).Decode(&deleteItems); err != nil {
		logger.Log.Debug("Cannot decode request JSON", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	delete := make([]model.URLToDelete, 0, len(deleteItems))
	for _, item := range deleteItems {
		delete = append(delete, model.URLToDelete{
			ShortLink: item,
			UserID:    userID.Value,
		})
	}

	s.service.ScheduleURLDeletion(r.Context(), delete)

	logger.Log.Debug("Urls deleted")
	w.WriteHeader(http.StatusAccepted)
}

// Stats - обрабатывает HTTP GET-запрос на получение статистики сервиса.
//
// Доступ к эндпоинту разрешён только если IP-адрес клиента из заголовка X-Real-IP
// принадлежит доверенной подсети, указанной в конфигурации (trusted_subnet).
//
// Если trusted_subnet не задан или IP не входит в подсеть — возвращает HTTP 403 Forbidden.
// В случае ошибки разбора CIDR — возвращает HTTP 500 Internal Server Error.
// В случае ошибки получения статистики — возвращает HTTP 500 Internal Server Error.
// При успешном выполнении возвращает HTTP 200 и JSON-объект с числом пользователей и URL.
func (s ShortenerHandler) Stats(w http.ResponseWriter, r *http.Request) {
	if s.config.TrustedSubnet == "" {
		logger.Log.Debug("Empty trusted subnet")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	_, IPNet, err := net.ParseCIDR(s.config.TrustedSubnet)
	if err != nil {
		logger.Log.Debug("Cannot parse CIDR", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	IP := net.ParseIP(r.Header.Get("X-Real-IP"))
	if IP == nil || !IPNet.Contains(IP) {
		logger.Log.Debug("Disallowed")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	stats, err := s.service.Stats(r.Context())
	if err != nil {
		logger.Log.Debug("Error getting stats", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusOK, stats)
}
