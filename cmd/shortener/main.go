package main

import (
	"net/http"

	"github.com/bubaew95/yandex-go-learn/internal/app"
	"github.com/go-chi/chi/v5"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	r := Routers()
	return http.ListenAndServe(`:8080`, r)
}

func Routers() chi.Router {
	urls := make(map[string]string)

	r := chi.NewRouter()
	r.Post("/", app.CreateURL(urls))
	r.Get("/{id}", app.GetURL(urls))

	return r
}
