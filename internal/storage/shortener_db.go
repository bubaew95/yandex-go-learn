package storage

import (
	"bufio"
	"encoding/json"
	"io"
	"os"

	"github.com/bubaew95/yandex-go-learn/config"
	"github.com/bubaew95/yandex-go-learn/internal/models"
)

type ShortenerDB struct {
	config   config.Config
	producer *Producer
}

func NewShortenerDB(c config.Config) (*ShortenerDB, error) {
	producer, err := NewProducer(c.FilePath)
	if err != nil {
		return nil, err
	}

	return &ShortenerDB{
		config:   c,
		producer: producer,
	}, nil
}

func (s ShortenerDB) Save(data *models.ShortenURL) error {
	err := s.producer.WriteShortener(data)
	if err != nil {
		return err
	}

	return nil
}

func (s ShortenerDB) Load() (map[string]string, error) {
	file, err := os.OpenFile(s.config.FilePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	bufio := bufio.NewReader(file)
	data := make(map[string]string)
	for {
		line, err := bufio.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		var s models.ShortenURL
		err = json.Unmarshal([]byte(line), &s)
		if err != nil {
			return nil, err
		}

		data[s.ShortURL] = s.OriginalURL
	}

	return data, nil
}

func (s ShortenerDB) Close() error {
	return s.producer.Close()
}
