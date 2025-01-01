package main

import (
	"fmt"
	"net/http"

	"github.com/bubaew95/yandex-go-learn/config"
	"github.com/bubaew95/yandex-go-learn/internal/app"
)

func main() {
	cfg := config.NewConfig()
	appInstance := app.NewApp(cfg)

	appInstance.Routers()

	if err := run(cfg, appInstance); err != nil {
		panic(err)
	}
}

func run(cfg *config.Config, app *app.App) error {
	fmt.Printf("Run server on port %s", cfg.Port)

	return http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), &app.Router)
}
