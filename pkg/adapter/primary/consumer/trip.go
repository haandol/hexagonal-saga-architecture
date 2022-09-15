package consumer

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/haandol/hexagonal/pkg/message"
	"github.com/haandol/hexagonal/pkg/message/event"
	"github.com/haandol/hexagonal/pkg/port/primaryport/consumerport"
	"github.com/haandol/hexagonal/pkg/service"
	"github.com/haandol/hexagonal/pkg/util"
)

type TripConsumer struct {
	*KafkaConsumer
	tripService *service.TripService
}

func NewTripConsumer(
	kafkaConsumer *KafkaConsumer,
	tripService *service.TripService,
) *TripConsumer {
	return &TripConsumer{
		KafkaConsumer: kafkaConsumer,
		tripService:   tripService,
	}
}

func (c *TripConsumer) Init() {
	logger := util.GetLogger().With(
		"module", "TripConsumer",
		"func", "Init",
	)

	if err := c.RegisterHandler(c.Handle); err != nil {
		logger.Panicw("Failed to register handler", "err", err.Error())
	}
}

func (c *TripConsumer) Handle(ctx context.Context, r *consumerport.Message) error {
	logger := util.GetLogger().With(
		"module", "TripConsumer",
		"func", "Handle",
	)

	msg := &message.Message{}
	if err := json.Unmarshal(r.Value, msg); err != nil {
		logger.Errorw("Failed to unmarshal command", "err", err.Error())
	}

	logger.Infow("Received command", "command", msg)
	con, seg := util.BeginSegmentWithTraceID(ctx, msg.CorrelationID, msg.ParentID, "## TripConsumer")
	seg.AddMetadata("msg", msg)
	defer seg.Close(nil)

	switch msg.Name {
	case "SagaEnded":
		evt := &event.SagaEnded{}
		if err := json.Unmarshal(r.Value, evt); err != nil {
			logger.Errorw("Failed to unmarshal command", "err", err.Error())
			seg.AddError(err)
			return err
		}
		return c.tripService.ProcessSagaEnded(con, evt)
	case "SagaAborted":
		evt := &event.SagaAborted{}
		if err := json.Unmarshal(r.Value, evt); err != nil {
			logger.Errorw("Failed to unmarshal command", "err", err.Error())
			seg.AddError(err)
			return err
		}
		return c.tripService.ProcessSagaAborted(con, evt)
	default:
		logger.Errorw("unknown command", "message", msg)
		err := errors.New("unknown command")
		seg.AddError(err)
		return err
	}
}
