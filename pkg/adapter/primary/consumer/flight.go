package consumer

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/haandol/hexagonal/pkg/message"
	"github.com/haandol/hexagonal/pkg/message/command"
	"github.com/haandol/hexagonal/pkg/port/primaryport/consumerport"
	"github.com/haandol/hexagonal/pkg/service"
	"github.com/haandol/hexagonal/pkg/util"
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
		logger.Panicw("Failed to register handler", "err", err.Error())
	}
}

func (c *FlightConsumer) Handle(ctx context.Context, r *consumerport.Message) error {
	logger := util.GetLogger().With(
		"module", "FlightConsumer",
		"func", "Handle",
	)

	msg := &message.Message{}
	if err := json.Unmarshal(r.Value, msg); err != nil {
		logger.Errorw("Failed to unmarshal command", "err", err.Error())
	}

	logger.Infow("Received command", "command", msg)
	con, seg := util.BeginSegmentWithTraceID(ctx, msg.CorrelationID, msg.ParentID, "## FlightConsumer")
	seg.AddMetadata("msg", msg)
	defer seg.Close(nil)

	switch msg.Name {
	case "BookFlight":
		cmd := &command.BookFlight{}
		if err := json.Unmarshal(r.Value, cmd); err != nil {
			logger.Errorw("Failed to unmarshal command", "err", err.Error())
			seg.AddError(err)
			return err
		}
		return c.flightService.Book(con, cmd)
	case "CancelFlightBooking":
		cmd := &command.CancelFlightBooking{}
		if err := json.Unmarshal(r.Value, cmd); err != nil {
			logger.Errorw("Failed to unmarshal command", "err", err.Error())
			seg.AddError(err)
			return err
		}
		return c.flightService.CancelBooking(con, cmd)
	default:
		logger.Errorw("unknown command", "message", msg)
		err := errors.New("unknown command")
		seg.AddError(err)
		return err
	}
}
