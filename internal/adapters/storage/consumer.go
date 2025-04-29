// Package storage предоставляет функции для работы с файловым хранилищем сокращённых ссылок.
package storage

import (
	"bufio"
	"encoding/json"
	"io"
	"os"

	"github.com/bubaew95/yandex-go-learn/internal/core/model"
)

// ReadShorteners читает файл с сериализованными записями сокращённых ссылок.
//
// Файл должен содержать строки в формате JSON, каждая строка — это структура model.ShortenURL.
// Каждая строка разбирается, и на её основе формируется карта соответствий: ключ — ShortURL, значение — OriginalURL.
//
// Если файл не существует, он будет создан с пустым содержимым.
//
// Возвращает карту сокращённых ссылок и ошибку, если операция чтения или десериализации завершилась неудачно.
func ReadShorteners(filename string) (map[string]string, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
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
