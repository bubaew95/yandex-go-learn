package app

import (
	"context"
	"errors"
	"github.com/bubaew95/yandex-go-learn/config"
	rpcServer "github.com/bubaew95/yandex-go-learn/internal/adapters/handlers/grpc"
	handlers "github.com/bubaew95/yandex-go-learn/internal/adapters/handlers/http"
	"github.com/bubaew95/yandex-go-learn/internal/adapters/handlers/http/middleware"
	"github.com/bubaew95/yandex-go-learn/internal/adapters/logger"
	"github.com/bubaew95/yandex-go-learn/internal/core/service"
	pb "github.com/bubaew95/yandex-go-learn/internal/proto"
	"github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"sync"
	"time"
)

type App struct {
	cfg                 *config.Config
	httpServer          *http.Server
	grpcServer          *grpc.Server
	shortenerRepository service.ShortenerRepository
	handler             handlers.ShortenerHandler
	wg                  sync.WaitGroup
}

func NewApp(cfg *config.Config, h handlers.ShortenerHandler, r service.ShortenerRepository) *App {
	return &App{
		cfg:                 cfg,
		shortenerRepository: r,
		handler:             h,
	}
}

func (a *App) Run() {
	if a.cfg.ListenHTTP {
		a.httpServer = listenHTTP(a.cfg, a.handler)
	}

	if a.cfg.ListenGRPC {
		a.grpcServer = listenGRPC(a.cfg, a.shortenerRepository)
	}
}

func (a *App) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if a.cfg.ListenHTTP && a.httpServer != nil {
		logger.Log.Info("Stopping HTTP server...")
		if err := a.httpServer.Shutdown(ctx); err != nil {
			logger.Log.Error("HTTP shutdown error", zap.Error(err))
		}
	}

	if a.cfg.ListenGRPC && a.grpcServer != nil {
		logger.Log.Info("Stopping gRPC server...")
		a.grpcServer.Stop()
	}
}

func setupRouter(shortenerHandler handlers.ShortenerHandler) *chi.Mux {
	route := chi.NewRouter()
	route.Use(middleware.LoggerMiddleware)
	route.Use(middleware.GZipMiddleware)
	route.Use(middleware.CookieMiddleware)

	route.Post("/", shortenerHandler.CreateURL)
	route.Get("/{id}", shortenerHandler.GetURL)
	route.Get("/ping", shortenerHandler.Ping)

	route.Route("/api/shorten", func(r chi.Router) {
		r.Post("/", shortenerHandler.AddNewURL)
		r.Post("/batch", shortenerHandler.Batch)
	})

	route.Route("/api/user", func(r chi.Router) {
		r.Get("/urls", shortenerHandler.GetUserURLS)
		r.Delete("/urls", shortenerHandler.DeleteUserURLS)
	})

	route.Route("/api/internal/stats", func(r chi.Router) {
		r.Get("/", shortenerHandler.Stats)
	})

	route.Mount("/debug", chi_middleware.Profiler())

	return route
}

func listenHTTP(cfg *config.Config, shortenerHandler handlers.ShortenerHandler) *http.Server {
	route := setupRouter(shortenerHandler)

	server := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: route,
	}

	if cfg.EnableHTTPS {
		logger.Log.Info("Running https server", zap.String("port", cfg.ServerAddress))
		go func() {
			if err := server.Serve(autocert.NewListener("example.com")); err != nil {
				logger.Log.Error("Failed to start https(tsl) server", zap.Error(err))
			}
		}()
	} else {
		logger.Log.Info("Running server", zap.String("ports", cfg.ServerAddress))
		go func() {
			if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				logger.Log.Error("Failed to start http server", zap.Error(err))
			}
		}()
	}

	return server
}

func listenGRPC(cfg *config.Config, sr service.ShortenerRepository) *grpc.Server {
	listener, err := net.Listen("tcp", ":3232")
	if err != nil {
		logger.Log.Fatal("Rpc server error", zap.Error(err))
	}

	shortenerService := service.NewShortenerService(sr, *cfg)

	s := grpc.NewServer(
		grpc.UnaryInterceptor(rpcServer.AuthInterceptor()),
	)
	pb.RegisterURLShortenerServer(s, rpcServer.NewServer(shortenerService))

	logger.Log.Info("Run rpc server")
	if err := s.Serve(listener); err != nil {
		logger.Log.Fatal("Rpc server error", zap.Error(err))
	}

	return s
}
