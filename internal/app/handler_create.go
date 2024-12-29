package app

import (
	"fmt"
	"io"
	"net/http"

	"github.com/bubaew95/yandex-go-learn/internal/utils"
)

func CreateURL(urls map[string]string) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		responseData, err := io.ReadAll(req.Body)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		body := string(responseData)
		res.WriteHeader(http.StatusCreated)
		res.Header().Set("content-type", "text/plain")
		res.Header().Set("content-length", string(len(responseData)))

		genID := utils.RandStringBytes(8)
		urls[genID] = body

		url := fmt.Sprintf("%s:%s/%s", "http://localhost", "8080", genID)
		res.Write([]byte(url))
	}
}
