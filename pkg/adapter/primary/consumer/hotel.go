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
	logger := util.GetLogger().With(
		"module", "HotelConsumer",
		"func", "Init",
	)

	if err := c.RegisterHandler(c.Handle); err != nil {
		logger.Panicw("Failed to register handler", "err", err.Error())
	}
}

func (c *HotelConsumer) Handle(ctx context.Context, r *consumerport.Message) error {
	logger := util.GetLogger().With(
		"module", "HotelConsumer",
		"func", "Handle",
	)

	msg := &message.Message{}
	if err := json.Unmarshal(r.Value, msg); err != nil {
		logger.Errorw("Failed to unmarshal command", "err", err.Error())
	}

	logger.Infow("Received command", "command", msg)
	con, seg := util.BeginSegmentWithTraceID(ctx, msg.CorrelationID, msg.ParentID, "## HotelConsumer")
	seg.AddMetadata("msg", msg)
	defer seg.Close(nil)

	switch msg.Name {
	case "BookHotel":
		cmd := &command.BookHotel{}
		if err := json.Unmarshal(r.Value, cmd); err != nil {
			logger.Errorw("Failed to unmarshal command", "err", err.Error())
			seg.AddError(err)
			return err
		}
		return c.hotelService.Book(con, cmd)
	case "CancelHotelBooking":
		cmd := &command.CancelHotelBooking{}
		if err := json.Unmarshal(r.Value, cmd); err != nil {
			logger.Errorw("Failed to unmarshal command", "err", err.Error())
			seg.AddError(err)
			return err
		}
		return c.hotelService.CancelBooking(con, cmd)
	default:
		logger.Errorw("unknown command", "message", msg)
		err := errors.New("unknown command")
		seg.AddError(err)
		return err
	}
}
