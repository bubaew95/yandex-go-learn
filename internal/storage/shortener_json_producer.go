package storage

import (
	"encoding/json"
	"os"

	"github.com/bubaew95/yandex-go-learn/internal/models"
)

type Producer struct {
	file    *os.File
	encoder *json.Encoder
}

func NewProducer(filename string) (*Producer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}

	return &Producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (p *Producer) WriteShortener(s *models.ShortenURL) error {
	return p.encoder.Encode(s)
}

func (p *Producer) Close() error {
	if err := p.file.Sync(); err != nil {
		return err
	}

	return p.file.Close()
}
