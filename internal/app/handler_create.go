package app

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/bubaew95/yandex-go-learn/internal/utils"
)

const randomStringLength = 8

func (app *App) CreateURL(res http.ResponseWriter, req *http.Request) {
	responseData, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	body := string(responseData)
	if body == "" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	genID := generateID(app.URLs)
	app.URLs[genID] = body
	url := fmt.Sprintf("%s/%s", app.Config.BaseURL, genID)

	res.WriteHeader(http.StatusCreated)
	res.Header().Set("content-type", "text/plain")
	res.Header().Set("content-length", strconv.Itoa(len(url)))

	res.Write([]byte(url))
}

func generateID(urls map[string]string) string {
	var genID string
	for {
		genID = utils.RandStringBytes(randomStringLength)
		if _, exists := urls[genID]; !exists {
			break
		}
	}

	return genID
}
