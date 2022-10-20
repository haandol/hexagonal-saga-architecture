package util

import (
	"context"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *Logger
var once sync.Once

type Logger struct {
	*zap.SugaredLogger
}

func InitLogger(stage string) *Logger {
	once.Do(func() {
		if strings.EqualFold(stage, "local") {
			cfg := zap.NewDevelopmentConfig()
			cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
			l, _ := cfg.Build()
			logger = &Logger{
				l.Sugar(),
			}
		} else {
			l, _ := zap.NewProduction()
			logger = &Logger{
				l.Sugar(),
			}
		}
	})
	logger.Infow("Logger initialized", "stage", stage)

	return logger
}

func GetLogger() *Logger {
	return logger
}

func (l *Logger) WithContext(ctx context.Context) *zap.SugaredLogger {
	segID := GetSegmentID(ctx)
	return logger.With(
		"TraceId", segID,
	)
}
