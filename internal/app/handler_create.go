package app

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"

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

	var mutex sync.Mutex
	ch := make(chan string)

	go generateID(app.URLs, &mutex, ch)

	genID := <-ch

	app.URLs[genID] = body
	url := fmt.Sprintf("%s/%s", app.Config.BaseURL, genID)

	res.WriteHeader(http.StatusCreated)
	res.Header().Set("content-type", "text/plain")
	res.Header().Set("content-length", strconv.Itoa(len(url)))

	res.Write([]byte(url))
}

func generateID(urls map[string]string, mutex *sync.Mutex, ch chan string) {
	var genID string
	for {
		genID = utils.RandStringBytes(randomStringLength)

		mutex.Lock()
		if _, exists := urls[genID]; !exists {
			urls[genID] = ""
			mutex.Unlock()
			break
		}
		mutex.Unlock()
	}
	ch <- genID
}
