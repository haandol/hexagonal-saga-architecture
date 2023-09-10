package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/haandol/hexagonal/pkg/constant"
	"github.com/haandol/hexagonal/pkg/message"
	"github.com/haandol/hexagonal/pkg/message/event"
	"github.com/haandol/hexagonal/pkg/port/primaryport/consumerport"
	"github.com/haandol/hexagonal/pkg/service"
	"github.com/haandol/hexagonal/pkg/util"
	"github.com/haandol/hexagonal/pkg/util/o11y"
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
	logger := util.GetLogger().WithGroup("TripConsumer.Init")

	if err := c.RegisterHandler(c.Handle); err != nil {
		logger.Error(constant.ErrTxtRegisterHandler, "err", err)
		panic(constant.ErrTxtRegisterHandler)
	}
}

func (c *TripConsumer) Handle(ctx context.Context, r *consumerport.Message) error {
	logger := util.GetLogger().WithGroup("TripConsumer.Handle")

	msg := &message.Message{}
	if err := json.Unmarshal(r.Value, msg); err != nil {
		logger.Error(constant.ErrTxtUnMarshalCommand, "err", err)
	}

	logger.Info("Received command", "command", msg)
	ctx, span := o11y.BeginSpanWithTraceID(ctx, msg.CorrelationID, msg.ParentID, "TripConsumer")
	defer span.End()
	span.SetAttributes(
		o11y.AttrString("msg", fmt.Sprintf("%v", msg)),
	)

	switch msg.Name {
	case "SagaEnded":
		evt := &event.SagaEnded{}
		if err := json.Unmarshal(r.Value, evt); err != nil {
			logger.Error(constant.ErrTxtUnMarshalCommand, "err", err)
			span.RecordError(err)
			span.SetStatus(o11y.GetStatus(err))
			span.SetStatus(o11y.GetStatus(err))
			return err
		}
		return c.tripService.ProcessSagaEnded(ctx, evt)
	case "SagaAborted":
		evt := &event.SagaAborted{}
		if err := json.Unmarshal(r.Value, evt); err != nil {
			logger.Error(constant.ErrTxtUnMarshalCommand, "err", err)
			span.RecordError(err)
			span.SetStatus(o11y.GetStatus(err))
			return err
		}
		return c.tripService.ProcessSagaAborted(ctx, evt)
	default:
		logger.Error(constant.ErrTxtUnknownCommand, "message", msg)
		err := errors.New(constant.ErrTxtUnknownCommand)
		span.RecordError(err)
		span.SetStatus(o11y.GetStatus(err))
		return err
	}
}
