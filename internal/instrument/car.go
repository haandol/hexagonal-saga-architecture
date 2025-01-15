package instrument

import (
	"log/slog"

	"github.com/haandol/hexagonal/pkg/o11y"
	"go.opentelemetry.io/otel/trace"
)

func RecordBookCarError(logger *slog.Logger, span trace.Span, err error, req interface{}) {
	logger.Error("failed to book car", "req", req, "err", err)
	span.RecordError(err)
	span.SetStatus(o11y.GetStatus(err))
}

func RecordCancelCarBookingError(logger *slog.Logger, span trace.Span, err error, cmd interface{}) {
	logger.Error("failed to cancel car booking", "err", err, "command", cmd)
	span.RecordError(err)
	span.SetStatus(o11y.GetStatus(err))
}

func RecordPublishCarBookedError(logger *slog.Logger, span trace.Span, err error, evt interface{}) {
	logger.Error("failed to publish CarBooked", "err", err, "evt", evt)
	span.RecordError(err)
	span.SetStatus(o11y.GetStatus(err))
}
