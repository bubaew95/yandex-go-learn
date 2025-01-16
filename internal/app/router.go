package app

import "github.com/bubaew95/yandex-go-learn/internal/service"

func (app *App) Routers() {
	app.Router.Use(service.LoggerMiddleware)
	app.Router.Post("/", app.CreateURL)
	app.Router.Get("/{id}", app.GetURL)
}
