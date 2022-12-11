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
	ctx, span := o11y.BeginSpanWithTraceID(ctx, msg.CorrelationID, msg.ParentID, "CarConsumer")
	defer span.End()
	span.SetAttributes(
		o11y.AttrString("msg", fmt.Sprintf("%v", msg)),
	)

	switch msg.Name {
	case "BookCar":
		cmd := &command.BookCar{}
		if err := json.Unmarshal(r.Value, cmd); err != nil {
			logger.Errorw("Failed to unmarshal command", "err", err.Error())
			span.RecordError(err)
			span.SetStatus(o11y.GetStatus(err))
			return err
		}
		return c.carService.Book(ctx, cmd)
	case "CancelCarBooking":
		cmd := &command.CancelCarBooking{}
		if err := json.Unmarshal(r.Value, cmd); err != nil {
			logger.Errorw("Failed to unmarshal command", "err", err.Error())
			span.RecordError(err)
			span.SetStatus(o11y.GetStatus(err))
			return err
		}
		return c.carService.CancelBooking(ctx, cmd)
	default:
		logger.Errorw("unknown command", "message", msg)
		err := errors.New("unknown command")
		span.RecordError(err)
		span.SetStatus(o11y.GetStatus(err))
		return err
	}
}
