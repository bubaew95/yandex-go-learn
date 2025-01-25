package storage

import (
	"os"

	"github.com/bubaew95/yandex-go-learn/config"
	"github.com/bubaew95/yandex-go-learn/internal/models"
)

type ShortenerDB struct {
	config config.Config
}

func NewShortenerDB(c config.Config) *ShortenerDB {
	return &ShortenerDB{
		config: c,
	}
}

func (s ShortenerDB) Save(data *models.ShortenURL) error {
	producer, err := NewProducer(s.config.FilePath)
	if err != nil {
		return err
	}
	defer producer.Close()

	err = producer.WriteShortener(data)
	if err != nil {
		return err
	}

	return nil
}

func (s ShortenerDB) Load() (map[string]string, error) {
	consumer, err := NewConsumer(s.config.FilePath)
	if err != nil {
		return nil, err
	}

	defer consumer.Close()

	read, err := consumer.ReadShorteners()
	if err != nil {
		return nil, err
	}

	return read, nil
}

func (s ShortenerDB) Count() int {
	datas, err := s.Load()
	if err != nil {
		return 0
	}

	return len(datas)
}

func (s ShortenerDB) RemoveFile() error {
	return os.Remove(s.config.FilePath)
}
