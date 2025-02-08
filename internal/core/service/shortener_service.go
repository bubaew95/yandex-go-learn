package service

import (
	"context"
	"fmt"
	"math/rand"
	"sync"

	"github.com/bubaew95/yandex-go-learn/config"
	"github.com/bubaew95/yandex-go-learn/internal/core/model"
	"github.com/bubaew95/yandex-go-learn/internal/core/ports"
)

type ShortenerService struct {
	repository ports.ShortenerRepositoryInterface
	config     config.Config
	mx         *sync.Mutex
}

func NewShortenerService(r ports.ShortenerRepositoryInterface, cfg config.Config) *ShortenerService {
	return &ShortenerService{
		repository: r,
		config:     cfg,
		mx:         &sync.Mutex{},
	}
}

func (s *ShortenerService) GenerateURL(ctx context.Context, url string, randomStringLength int) string {
	s.mx.Lock()
	defer s.mx.Unlock()

	var genID string
	for {
		genID = s.RandStringBytes(randomStringLength)

		_, existsURL := s.repository.GetURLByID(ctx, genID)
		if !existsURL {
			s.repository.SetURL(ctx, genID, url)
			break
		}
	}

	return s.generateResponseURL(genID)
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
	s.mx.Lock()
	defer s.mx.Unlock()

	return s.repository.InsertURLs(ctx, urls)
}
