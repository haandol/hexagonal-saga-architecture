package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/haandol/hexagonal/internal/constant"
	"github.com/haandol/hexagonal/internal/message"
	"github.com/haandol/hexagonal/internal/message/command"
	"github.com/haandol/hexagonal/internal/port/primaryport/consumerport"
	"github.com/haandol/hexagonal/internal/service"
	"github.com/haandol/hexagonal/pkg/o11y"
	"github.com/haandol/hexagonal/pkg/util"
	"github.com/pkg/errors"
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
	logger := util.GetLogger().WithGroup("CarConsumer.Init")

	if err := c.RegisterHandler(c.Handle); err != nil {
		logger.Error(constant.ErrTxtRegisterHandler, "err", err)
		panic(constant.ErrTxtRegisterHandler)
	}
}

func (c *CarConsumer) Handle(ctx context.Context, r *consumerport.Message) error {
	logger := util.GetLogger().WithGroup("CarConsumer.Handle")

	msg := &message.Message{}
	if err := json.Unmarshal(r.Value, msg); err != nil {
		logger.Error(constant.ErrTxtUnMarshalCommand, "err", err)
	}

	logger.Info("Received command", "command", msg)
	ctx, span := o11y.BeginSpanWithTraceID(ctx, msg.CorrelationID, msg.ParentID, "CarConsumer")
	defer span.End()
	span.SetAttributes(
		o11y.AttrString("msg", fmt.Sprintf("%v", msg)),
	)

	switch msg.Name {
	case "BookCar":
		cmd := &command.BookCar{}
		if err := json.Unmarshal(r.Value, cmd); err != nil {
			logger.Error(constant.ErrTxtUnMarshalCommand, "err", err)
			span.RecordError(err)
			span.SetStatus(o11y.GetStatus(err))
			return err
		}
		return c.carService.Book(ctx, cmd)
	case "CancelCarBooking":
		cmd := &command.CancelCarBooking{}
		if err := json.Unmarshal(r.Value, cmd); err != nil {
			logger.Error(constant.ErrTxtUnMarshalCommand, "err", err)
			span.RecordError(err)
			span.SetStatus(o11y.GetStatus(err))
			return err
		}
		return c.carService.CancelBooking(ctx, cmd)
	default:
		logger.Error(constant.ErrTxtUnknownCommand, "message", msg)
		err := errors.New(constant.ErrTxtUnknownCommand)
		span.RecordError(err)
		span.SetStatus(o11y.GetStatus(err))
		return err
	}
}
