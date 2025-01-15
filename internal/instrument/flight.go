package instrument

import (
	"log/slog"

	"github.com/haandol/hexagonal/pkg/o11y"
	"go.opentelemetry.io/otel/trace"
)

func RecordBookFlightError(logger *slog.Logger, span trace.Span, err error, req interface{}) {
	logger.Error("failed to book flight", "req", req, "err", err)
	span.RecordError(err)
	span.SetStatus(o11y.GetStatus(err))
}

func RecordCancelFlightBookingError(logger *slog.Logger, span trace.Span, err error, cmd interface{}) {
	logger.Error("failed to cancel flight booking", "err", err, "command", cmd)
	span.RecordError(err)
	span.SetStatus(o11y.GetStatus(err))
}
