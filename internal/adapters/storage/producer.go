// Package storage содержит утилиты для сериализации и десериализации сокращённых URL в файловом формате.
package storage

import (
	"encoding/json"
	"os"

	"github.com/bubaew95/yandex-go-learn/internal/core/model"
)

// Producer реализует механизм последовательной записи JSON-записей в файл.
// Используется для сохранения сокращённых ссылок в формате model.ShortenURL.
type Producer struct {
	file    *os.File
	encoder *json.Encoder
}

// NewProducer открывает (или создаёт) файл по указанному пути и возвращает новый экземпляр Producer.
//
// Файл открывается в режиме дозаписи (append), таким образом новые записи не затирают старые.
// Права доступа к файлу устанавливаются как 0777.
//
// Возвращает ошибку при неудачном открытии файла.
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

// WriteShortener сериализует структуру model.ShortenURL и записывает её в файл в формате JSON.
//
// Каждая запись пишется как отдельная строка.
func (p *Producer) WriteShortener(s *model.ShortenURL) error {
	return p.encoder.Encode(s)
}

// Close завершает работу с файлом: вызывает синхронизацию буфера и закрывает файл.
//
// Возвращает ошибку, если одна из операций завершилась неудачно.
func (p *Producer) Close() error {
	if err := p.file.Sync(); err != nil {
		return err
	}

	return p.file.Close()
}
