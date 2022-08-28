package util

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *Logger
var once sync.Once

type Logger struct {
	*zap.SugaredLogger
}

func GetLogger() *Logger {
	once.Do(func() {
		if "production" == os.Getenv("STAGE") {
			l, _ := zap.NewProduction()
			logger = &Logger{
				l.Sugar(),
			}
		} else {
			cfg := zap.NewDevelopmentConfig()
			cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
			l, _ := cfg.Build()
			logger = &Logger{
				l.Sugar(),
			}
		}
	})
	return logger
}
