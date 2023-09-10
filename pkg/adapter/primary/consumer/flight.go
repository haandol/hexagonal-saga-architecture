package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/haandol/hexagonal/pkg/message"
	"github.com/haandol/hexagonal/pkg/message/command"
	"github.com/haandol/hexagonal/pkg/port/primaryport/consumerport"
	"github.com/haandol/hexagonal/pkg/service"
	"github.com/haandol/hexagonal/pkg/util"
	"github.com/haandol/hexagonal/pkg/util/o11y"
)

type FlightConsumer struct {
	*KafkaConsumer
	flightService *service.FlightService
}

func NewFlightConsumer(
	kafkaConsumer *KafkaConsumer,
	flightService *service.FlightService,
) *FlightConsumer {
	return &FlightConsumer{
		KafkaConsumer: kafkaConsumer,
		flightService: flightService,
	}
}

func (c *FlightConsumer) Init() {
	logger := util.GetLogger().With(
		"module", "FlightConsumer",
		"func", "Init",
	)

	if err := c.RegisterHandler(c.Handle); err != nil {
		msg := "Failed to register handler"
		logger.Error(msg, "err", err)
		panic(msg)
	}
}

func (c *FlightConsumer) Handle(ctx context.Context, r *consumerport.Message) error {
	logger := util.GetLogger().With(
		"module", "FlightConsumer",
		"func", "Handle",
	)

	msg := &message.Message{}
	if err := json.Unmarshal(r.Value, msg); err != nil {
		logger.Error("Failed to unmarshal command", "err", err)
	}

	logger.Info("Received command", "command", msg)
	ctx, span := o11y.BeginSpanWithTraceID(ctx, msg.CorrelationID, msg.ParentID, "FlightConsumer")
	defer span.End()
	span.SetAttributes(
		o11y.AttrString("msg", fmt.Sprintf("%v", msg)),
	)

	switch msg.Name {
	case "BookFlight":
		cmd := &command.BookFlight{}
		if err := json.Unmarshal(r.Value, cmd); err != nil {
			logger.Error("Failed to unmarshal command", "err", err)
			span.RecordError(err)
			span.SetStatus(o11y.GetStatus(err))
			return err
		}
		return c.flightService.Book(ctx, cmd)
	case "CancelFlightBooking":
		cmd := &command.CancelFlightBooking{}
		if err := json.Unmarshal(r.Value, cmd); err != nil {
			logger.Error("Failed to unmarshal command", "err", err)
			span.RecordError(err)
			span.SetStatus(o11y.GetStatus(err))
			return err
		}
		return c.flightService.CancelBooking(ctx, cmd)
	default:
		logger.Error("unknown command", "message", msg)
		err := errors.New("unknown command")
		span.RecordError(err)
		span.SetStatus(o11y.GetStatus(err))
		return err
	}
}
