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
	r := app.Routers()
	return http.ListenAndServe(`:8080`, r)
}
