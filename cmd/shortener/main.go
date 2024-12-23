package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
)

const PORT = "8080"
const DOMAIN = "http://localhost"
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	urls := make(map[string]string)

	mux := http.NewServeMux()
	mux.HandleFunc(`/`, func(res http.ResponseWriter, req *http.Request) {
		if http.MethodPost != req.Method {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		responseData, err := io.ReadAll(req.Body)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		body := string(responseData)

		res.WriteHeader(http.StatusCreated)
		res.Header().Set("content-type", "text/plain")
		res.Header().Set("content-length", "30")

		genID := RandStringBytes(8)
		urls[genID] = body

		url := fmt.Sprintf("%s:%s/%s", DOMAIN, PORT, genID)
		res.Write([]byte(url))
	})

	mux.HandleFunc(`/{id}`, func(res http.ResponseWriter, req *http.Request) {
		if http.MethodGet != req.Method {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		id := req.PathValue(`id`)

		url, ok := urls[id]
		if !ok {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		res.Header().Set("Location", url)
		res.WriteHeader(http.StatusTemporaryRedirect)
	})

	return http.ListenAndServe(`:`+PORT, mux)
}

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
