package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/haandol/hexagonal/internal/constant"
	"github.com/haandol/hexagonal/internal/message"
	"github.com/haandol/hexagonal/internal/message/command"
	"github.com/haandol/hexagonal/internal/port/primaryport/consumerport"
	"github.com/haandol/hexagonal/internal/service"
	"github.com/haandol/hexagonal/pkg/util"
	"github.com/haandol/hexagonal/pkg/util/o11y"
)

type HotelConsumer struct {
	*KafkaConsumer
	hotelService *service.HotelService
}

func NewHotelConsumer(
	kafkaConsumer *KafkaConsumer,
	hotelService *service.HotelService,
) *HotelConsumer {
	return &HotelConsumer{
		KafkaConsumer: kafkaConsumer,
		hotelService:  hotelService,
	}
}

func (c *HotelConsumer) Init() {
	logger := util.GetLogger().WithGroup("HotelConsumer.Init")

	if err := c.RegisterHandler(c.Handle); err != nil {
		logger.Error(constant.ErrTxtRegisterHandler, "err", err)
		panic(constant.ErrTxtRegisterHandler)
	}
}

func (c *HotelConsumer) Handle(ctx context.Context, r *consumerport.Message) error {
	logger := util.GetLogger().WithGroup("HotelConsumer.Handle")

	msg := &message.Message{}
	if err := json.Unmarshal(r.Value, msg); err != nil {
		logger.Error(constant.ErrTxtUnMarshalCommand, "err", err)
	}

	logger.Info("Received command", "command", msg)
	ctx, span := o11y.BeginSpanWithTraceID(ctx, msg.CorrelationID, msg.ParentID, "HotelConsumer")
	defer span.End()
	span.SetAttributes(
		o11y.AttrString("msg", fmt.Sprintf("%v", msg)),
	)

	switch msg.Name {
	case "BookHotel":
		cmd := &command.BookHotel{}
		if err := json.Unmarshal(r.Value, cmd); err != nil {
			logger.Error(constant.ErrTxtUnMarshalCommand, "err", err)
			span.RecordError(err)
			span.SetStatus(o11y.GetStatus(err))
			return err
		}
		return c.hotelService.Book(ctx, cmd)
	case "CancelHotelBooking":
		cmd := &command.CancelHotelBooking{}
		if err := json.Unmarshal(r.Value, cmd); err != nil {
			logger.Error(constant.ErrTxtUnMarshalCommand, "err", err)
			span.RecordError(err)
			span.SetStatus(o11y.GetStatus(err))
			return err
		}
		return c.hotelService.CancelBooking(ctx, cmd)
	default:
		logger.Error(constant.ErrTxtUnknownCommand, "message", msg)
		err := errors.New(constant.ErrTxtUnknownCommand)
		span.RecordError(err)
		span.SetStatus(o11y.GetStatus(err))
		return err
	}
}
