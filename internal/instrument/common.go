package instrument

import (
	"log/slog"

	"github.com/haandol/hexagonal/pkg/o11y"
	"go.opentelemetry.io/otel/trace"
)

func RecordBeginTxError(logger *slog.Logger, span trace.Span, err error) {
	logger.Error("failed to book car", "err", err)
	span.RecordError(err)
	span.SetStatus(o11y.GetStatus(err))
}

func RecordCommitTxError(logger *slog.Logger, span trace.Span, err error) {
	logger.Error("failed to commit transaction", "err", err)
	span.RecordError(err)
	span.SetStatus(o11y.GetStatus(err))
}

func RecordRollbackTxError(logger *slog.Logger, span trace.Span, err error) {
	logger.Error("failed to rollback transaction", "err", err)
	span.RecordError(err)
	span.SetStatus(o11y.GetStatus(err))
}

func RecordPublishAbortSagaError(logger *slog.Logger, span trace.Span, err error, cmd interface{}) {
	logger.Error("failed to publish AbortSaga command", "err", err, "cmd", cmd)
	span.RecordError(err)
	span.SetStatus(o11y.GetStatus(err))
}
