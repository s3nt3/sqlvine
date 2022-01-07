package logger

import (
	"go.uber.org/zap"
)

var (
	// global logger
	L *zap.SugaredLogger
)

func init() {
	// cfg := zap.NewProductionConfig()
	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{
		"/tmp/sqlvine/sqlvine.log",
	}
	logger, _ := cfg.Build()
	logger.Sync()
	L = logger.Sugar()
}
