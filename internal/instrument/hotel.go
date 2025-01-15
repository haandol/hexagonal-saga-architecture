package instrument

import (
	"log/slog"

	"github.com/haandol/hexagonal/pkg/o11y"
	"go.opentelemetry.io/otel/trace"
)

func RecordBookHotelError(logger *slog.Logger, span trace.Span, err error, req interface{}) {
	logger.Error("failed to book hotel", "req", req, "err", err)
	span.RecordError(err)
	span.SetStatus(o11y.GetStatus(err))
}

func RecordCancelHotelBookingError(logger *slog.Logger, span trace.Span, err error, cmd interface{}) {
	logger.Error("failed to cancel hotel booking", "err", err, "command", cmd)
	span.RecordError(err)
	span.SetStatus(o11y.GetStatus(err))
}
