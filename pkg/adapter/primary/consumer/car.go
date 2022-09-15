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

type CarConsumer struct {
	*KafkaConsumer
	carService *service.CarService
}

func NewCarConsumer(
	kafkaConsumer *KafkaConsumer,
	carService *service.CarService,
) *CarConsumer {
	return &CarConsumer{
		KafkaConsumer: kafkaConsumer,
		carService:    carService,
	}
}

func (c *CarConsumer) Init() {
	logger := util.GetLogger().With(
		"module", "CarConsumer",
		"func", "Init",
	)

	if err := c.RegisterHandler(c.Handle); err != nil {
		logger.Panicw("Failed to register handler", "err", err.Error())
	}
}

func (c *CarConsumer) Handle(ctx context.Context, r *consumerport.Message) error {
	logger := util.GetLogger().With(
		"module", "CarConsumer",
		"func", "Handle",
	)

	msg := &message.Message{}
	if err := json.Unmarshal(r.Value, msg); err != nil {
		logger.Errorw("Failed to unmarshal command", "err", err.Error())
	}

	logger.Infow("Received command", "command", msg)
	con, seg := util.BeginSegmentWithTraceID(ctx, msg.CorrelationID, "## CarConsumer")
	seg.AddMetadata("msg", msg)
	defer seg.Close(nil)

	switch msg.Name {
	case "BookCar":
		cmd := &command.BookCar{}
		if err := json.Unmarshal(r.Value, cmd); err != nil {
			logger.Errorw("Failed to unmarshal command", "err", err.Error())
			seg.AddError(err)
			return err
		}
		return c.carService.Book(con, cmd)
	case "CancelCarBooking":
		cmd := &command.CancelCarBooking{}
		if err := json.Unmarshal(r.Value, cmd); err != nil {
			logger.Errorw("Failed to unmarshal command", "err", err.Error())
			seg.AddError(err)
			return err
		}
		return c.carService.CancelBooking(con, cmd)
	default:
		logger.Errorw("unknown command", "message", msg)
		err := errors.New("unknown command")
		seg.AddError(err)
		return err
	}
}
