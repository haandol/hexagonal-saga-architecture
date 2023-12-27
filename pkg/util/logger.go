package util

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/haandol/hexagonal/pkg/o11y"
)

var logger *slog.Logger

func InitLogger(stage string) *slog.Logger {
	var once sync.Once
	once.Do(func() {
		if strings.EqualFold(stage, "local") {
			logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
		} else {
			logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
		}
	})
	logger.Info("Logger initialized", "stage", stage)

	return logger
}

func GetLogger() *slog.Logger {
	return logger
}

func LoggerFromContext(ctx context.Context) *slog.Logger {
	traceID, spanID := o11y.GetTraceSpanID(ctx)
	return logger.With(
		"TraceId", o11y.GetXrayTraceID(traceID),
		"SpanId", spanID,
	)
}
