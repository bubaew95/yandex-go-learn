package service

import (
	"sync"
)

type Storage struct {
	URLs map[string]string
	mx   *sync.Mutex
}

func NewStorage(URLs map[string]string) *Storage {
	return &Storage{
		URLs: URLs,
		mx:   &sync.Mutex{},
	}
}

func (s *Storage) GenerateID(url string, randomStringLength int) string {
	s.mx.Lock()
	defer s.mx.Unlock()

	var genID string
	for {
		genID = RandStringBytes(randomStringLength)

		_, exists := s.URLs[genID]
		if !exists {
			s.URLs[genID] = url
			break
		}
	}

	return genID
}
