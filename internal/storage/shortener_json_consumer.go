package storage

import (
	"bufio"
	"encoding/json"
	"io"
	"os"

	"github.com/bubaew95/yandex-go-learn/internal/models"
)

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

		var s models.ShortenURL
		err = json.Unmarshal([]byte(line), &s)
		if err != nil {
			return nil, err
		}

		data[s.ShortURL] = s.OriginalURL
	}

	return data, nil
}
