package app

func (app *App) Routers() {
	app.Router.Post("/", app.CreateURL)
	app.Router.Get("/{id}", app.GetURL)
}
