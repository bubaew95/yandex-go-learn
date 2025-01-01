package app

import (
	"fmt"
	"io"
	"net/http"

	"github.com/bubaew95/yandex-go-learn/internal/utils"
)

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

	res.WriteHeader(http.StatusCreated)
	res.Header().Set("content-type", "text/plain")
	res.Header().Set("content-length", "30")

	genID := utils.RandStringBytes(8)
	app.URLs[genID] = body

	url := fmt.Sprintf("%s/%s", app.Config.BaseURL, genID)
	res.Write([]byte(url))
}
