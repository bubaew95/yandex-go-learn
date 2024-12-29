package app

import (
	"github.com/go-chi/chi/v5"
)

func Routers() chi.Router {
	urls := make(map[string]string)

	r := chi.NewRouter()
	r.Post("/", CreateURL(urls))
	r.Get("/{id}", GetURL(urls))

	return r
}
