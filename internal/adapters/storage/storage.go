// Package storage предоставляет реализацию хранилища сокращённых URL,
// использующего файловую систему для записи и чтения данных в формате JSON.
package storage

import (
	"bufio"
	"encoding/json"
	"io"
	"os"

	"github.com/bubaew95/yandex-go-learn/config"
	"github.com/bubaew95/yandex-go-learn/internal/core/model"
)

// ShortenerDB реализует файловое хранилище сокращённых ссылок.
// Хранение данных осуществляется в виде последовательных JSON-записей.
type ShortenerDB struct {
	config   config.Config
	producer *Producer
}

// NewShortenerDB инициализирует файловое хранилище и готовит его к записи новых записей.
//
// Открывает файл, указанный в конфигурации, для последующей записи.
// Возвращает ошибку, если файл не удалось открыть.
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

// Save сериализует объект model.ShortenURL и добавляет его в файл хранилища.
//
// Возвращает ошибку, если не удалось записать данные.
func (s ShortenerDB) Save(data *model.ShortenURL) error {
	err := s.producer.WriteShortener(data)
	if err != nil {
		return err
	}

	return nil
}

// Load загружает все записи из файла и возвращает отображение ID -> оригинальный URL.
//
// Каждая строка файла должна быть JSON-представлением структуры model.ShortenURL.
// Возвращает ошибку при невозможности прочитать или десериализовать строку.
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

		var s model.ShortenURL
		err = json.Unmarshal([]byte(line), &s)
		if err != nil {
			return nil, err
		}

		data[s.ShortURL] = s.OriginalURL
	}

	return data, nil
}

// Close завершает работу с хранилищем, закрывая файловый поток записи.
//
// Возвращает ошибку, если операция завершения не удалась.
func (s ShortenerDB) Close() error {
	return s.producer.Close()
}
