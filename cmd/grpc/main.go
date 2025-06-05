package main

import (
	"fmt"
	"github.com/bubaew95/yandex-go-learn/config"
	rpcServer "github.com/bubaew95/yandex-go-learn/internal/adapters/handlers/grpc"
	"github.com/bubaew95/yandex-go-learn/internal/adapters/logger"
	"github.com/bubaew95/yandex-go-learn/internal/core/service"
	pb "github.com/bubaew95/yandex-go-learn/internal/proto"
	"github.com/bubaew95/yandex-go-learn/pkg/helper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

func main() {
	if err := logger.Initialize(); err != nil {
		fmt.Errorf("logging initialization error: %w", err)
	}

	listener, err := net.Listen("tcp", ":3232")
	if err != nil {
		logger.Log.Fatal("Rpc server error", zap.Error(err))
	}

	cfg := config.NewConfig()
	shortenerRepository, err := helper.InitRepositoryHelper(*cfg)
	if err != nil {
		logger.Log.Fatal("Repository init error", zap.Error(err))
	}
	defer helper.SafeClose(shortenerRepository)

	shortenerService := service.NewShortenerService(shortenerRepository, *cfg)

	s := grpc.NewServer(
		grpc.UnaryInterceptor(rpcServer.AuthInterceptor()),
	)
	pb.RegisterURLShortenerServer(s, rpcServer.NewServer(shortenerService))

	logger.Log.Info("Run rpc server")
	if err := s.Serve(listener); err != nil {
		logger.Log.Fatal("Rpc server error", zap.Error(err))
	}
}
