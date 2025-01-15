package instrument

import (
	"log/slog"

	"github.com/haandol/hexagonal/pkg/o11y"
	"go.opentelemetry.io/otel/trace"
)

func RecordStartSagaError(logger *slog.Logger, span trace.Span, err error, cmd interface{}) {
	logger.Error("failed to start saga", "err", err, "cmd", cmd)
	span.RecordError(err)
	span.SetStatus(o11y.GetStatus(err))
}

func RecordUpdateSagaError(logger *slog.Logger, span trace.Span, err error, tripID uint) {
	logger.Error("failed to update saga", "err", err, "tripID", tripID)
	span.RecordError(err)
	span.SetStatus(o11y.GetStatus(err))
}

func RecordAbortSagaError(logger *slog.Logger, span trace.Span, err error, cmd interface{}) {
	logger.Error("failed to abort saga", "err", err, "cmd", cmd)
	span.RecordError(err)
	span.SetStatus(o11y.GetStatus(err))
}

func RecordEndSagaError(logger *slog.Logger, span trace.Span, err error, cmd interface{}) {
	logger.Error("failed to end saga", "err", err, "cmd", cmd)
	span.RecordError(err)
	span.SetStatus(o11y.GetStatus(err))
}

func RecordPublishSagaCommandError(logger *slog.Logger, span trace.Span, err error, cmd interface{}) {
	logger.Error("failed to publish saga command", "err", err, "cmd", cmd)
	span.RecordError(err)
	span.SetStatus(o11y.GetStatus(err))
}

func RecordProcessSagaEventError(logger *slog.Logger, span trace.Span, err error, evt interface{}) {
	logger.Error("failed to process saga evt", "err", err, "evt", evt)
	span.RecordError(err)
	span.SetStatus(o11y.GetStatus(err))
}

func RecordCompensateSagaEventError(logger *slog.Logger, span trace.Span, err error, evt interface{}) {
	logger.Error("failed to compensate saga evt", "err", err, "evt", evt)
	span.RecordError(err)
	span.SetStatus(o11y.GetStatus(err))
}
