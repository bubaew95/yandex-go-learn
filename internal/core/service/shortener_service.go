package service

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/bubaew95/yandex-go-learn/config"
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

func (s *ShortenerService) GenerateURL(url string, randomStringLength int) string {
	s.mx.Lock()
	defer s.mx.Unlock()

	var genID string
	for {
		genID = s.RandStringBytes(randomStringLength)

		_, existsURL := s.repository.GetURLByID(genID)
		if !existsURL {
			s.repository.SetURL(genID, url)
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

func (s ShortenerService) GetURLByID(id string) (string, bool) {
	return s.repository.GetURLByID(id)
}

func (s ShortenerService) GetAllURL() map[string]string {
	return s.repository.GetAllURL()
}

func (s ShortenerService) Ping() error {
	return s.repository.Ping()
}

func (s ShortenerService) generateResponseURL(id string) string {
	return fmt.Sprintf("%s/%s", s.config.BaseURL, id)
}
