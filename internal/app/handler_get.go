package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *App) GetURL(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, `id`)

	url, ok := app.URLs[id]
	if !ok {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	res.Header().Set("Location", url)
	res.WriteHeader(http.StatusTemporaryRedirect)
}
