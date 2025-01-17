package service

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/bubaew95/yandex-go-learn/internal/interfaces"
)

type ShortenerService struct {
	repository interfaces.ShortenerRepositoryInterface
	mx         *sync.Mutex
}

func NewShortenerService(r interfaces.ShortenerRepositoryInterface) *ShortenerService {
	return &ShortenerService{
		repository: r,
		mx:         &sync.Mutex{},
	}
}

func (s *ShortenerService) GenerateID(url string, randomStringLength int) string {
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

	return genID
}

func (s ShortenerService) RandStringBytes(n int) string {

	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}

func (s ShortenerService) GenerateResponseURL(id string) string {
	return fmt.Sprintf("%s/%s", s.repository.GetBaseURL(), id)
}

func (s ShortenerService) GetURLByID(id string) (string, bool) {
	return s.repository.GetURLByID(id)
}
