package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log - логгер.
var Log *zap.Logger = zap.NewNop()

// Initialize инициализация логгера.
func Initialize() error {
	cfg := zap.NewProductionConfig()

	cfg.Level.SetLevel(zapcore.DebugLevel)

	zl, err := cfg.Build()
	if err != nil {
		return err
	}

	Log = zl
	return nil
}
