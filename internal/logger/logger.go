package logger

import "go.uber.org/zap"

var Log *zap.Logger = zap.NewNop()

func Initialize() error {
	cfg := zap.NewProductionConfig()

	zl, err := cfg.Build()
	if err != nil {
		return err
	}

	Log = zl
	return nil
}
