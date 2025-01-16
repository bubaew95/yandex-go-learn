package main

import (
	"fmt"
	"net/http"

	"github.com/bubaew95/yandex-go-learn/config"
	"github.com/bubaew95/yandex-go-learn/internal/app"
	"github.com/bubaew95/yandex-go-learn/internal/logger"
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
	if err := logger.Initialize(); err != nil {
		return err
	}

	fmt.Printf("Run server on port %s", cfg.Port)

	return http.ListenAndServe(cfg.Port, &app.Router)
}
