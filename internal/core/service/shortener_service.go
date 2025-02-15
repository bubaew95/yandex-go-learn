package service

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"sync"

	"github.com/bubaew95/yandex-go-learn/config"
	"github.com/bubaew95/yandex-go-learn/internal/core/model"
	"github.com/bubaew95/yandex-go-learn/internal/core/ports"
)

type ShortenerService struct {
	repository ports.ShortenerRepository
	config     config.Config
	mx         *sync.Mutex
}

func NewShortenerService(r ports.ShortenerRepository, cfg config.Config) *ShortenerService {
	return &ShortenerService{
		repository: r,
		config:     cfg,
		mx:         &sync.Mutex{},
	}
}

func (s ShortenerService) GenerateURL(ctx context.Context, url string, randomStringLength int) (string, error) {
	s.mx.Lock()
	defer s.mx.Unlock()

	var genID string
	for {
		genID = s.RandStringBytes(randomStringLength)

		_, existsURL := s.repository.GetURLByID(ctx, genID)
		if !existsURL {
			err := s.repository.SetURL(ctx, genID, url)
			if err != nil {
				return "", err
			}

			break
		}
	}

	return s.generateResponseURL(genID), nil
}

func (s ShortenerService) RandStringBytes(n int) string {

	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}

func (s ShortenerService) GetURLByID(ctx context.Context, id string) (string, bool) {
	return s.repository.GetURLByID(ctx, id)
}

func (s ShortenerService) GetURLByOriginalURL(ctx context.Context, originalURL string) (string, bool) {
	id, ok := s.repository.GetURLByOriginalURL(ctx, originalURL)

	if ok {
		return s.generateResponseURL(id), ok
	}

	return id, ok
}

func (s ShortenerService) GetAllURL(ctx context.Context) map[string]string {
	return s.repository.GetAllURL(ctx)
}

func (s ShortenerService) Ping() error {
	return s.repository.Ping()
}

func (s ShortenerService) generateResponseURL(id string) string {
	return fmt.Sprintf("%s/%s", s.config.BaseURL, id)
}

func (s ShortenerService) InsertURLs(ctx context.Context, urls []model.ShortenerURLMapping) ([]model.ShortenerURLResponse, error) {
	var items []model.ShortenerURLMapping
	for _, v := range urls {
		if isEmpty(v.CorrelationID) || isEmpty(v.OriginalURL) {
			continue
		}

		items = append(items, v)
	}

	err := s.repository.InsertURLs(ctx, items)
	if err != nil {
		return nil, err
	}

	var responseURLs []model.ShortenerURLResponse
	for _, v := range items {
		responseURLs = append(responseURLs, model.ShortenerURLResponse{
			CorrelationID: v.CorrelationID,
			ShortURL:      s.generateResponseURL(v.CorrelationID),
		})
	}

	return responseURLs, nil
}

func isEmpty(t string) bool {
	return strings.TrimSpace(t) == ""
}
