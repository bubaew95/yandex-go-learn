package main

import (
	"net/http"

	"github.com/bubaew95/yandex-go-learn/internal/app"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	urls := make(map[string]string)

	mux := http.NewServeMux()
	mux.HandleFunc(`/`, app.CreateURL(urls))
	mux.HandleFunc(`/{id}`, app.GetURL(urls))

	return http.ListenAndServe(`:8080`, mux)
}
