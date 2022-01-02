package logger

import (
	"go.uber.org/zap"
)

var (
	// global logger
	L *zap.SugaredLogger
)

func init() {
	// logger, _ := zap.NewProduction()
	logger, _ := zap.NewDevelopment()
	logger.Sync()
	L = logger.Sugar()
}
