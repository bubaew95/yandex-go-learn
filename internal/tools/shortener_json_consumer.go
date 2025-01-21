package tools

import (
	"bufio"
	"encoding/json"
	"io"
	"os"

	"github.com/bubaew95/yandex-go-learn/internal/models"
)

type Consumer struct {
	file  *os.File
	bufio *bufio.Reader
}

func NewConsumer(filename string) (*Consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		file:  file,
		bufio: bufio.NewReader(file),
	}, nil
}

func (s *Consumer) ReadShorteners() (map[string]string, error) {
	data := make(map[string]string)

	for {
		line, err := s.bufio.ReadString('\n')
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

		data[s.ShortUrl] = s.OriginalUrl
	}

	return data, nil
}

func (c *Consumer) Close() error {
	return c.file.Close()
}
