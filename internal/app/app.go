package app

import (
	"github.com/bubaew95/yandex-go-learn/config"
	"github.com/go-chi/chi/v5"
)

type App struct {
	URLs   map[string]string
	Config *config.Config
	Router chi.Mux
}

func NewApp(cfg *config.Config) *App {
	return &App{
		URLs:   make(map[string]string),
		Config: cfg,
		Router: *chi.NewRouter(),
	}
}
