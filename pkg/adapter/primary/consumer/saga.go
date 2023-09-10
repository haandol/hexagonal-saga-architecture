package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/haandol/hexagonal/pkg/constant"
	"github.com/haandol/hexagonal/pkg/message"
	"github.com/haandol/hexagonal/pkg/message/command"
	"github.com/haandol/hexagonal/pkg/message/event"
	"github.com/haandol/hexagonal/pkg/port/primaryport/consumerport"
	"github.com/haandol/hexagonal/pkg/service"
	"github.com/haandol/hexagonal/pkg/util"
	"github.com/haandol/hexagonal/pkg/util/o11y"
)

type SagaConsumer struct {
	*KafkaConsumer
	sagaService *service.SagaService
}

func NewSagaConsumer(
	kafkaConsumer *KafkaConsumer,
	sagaService *service.SagaService,
) *SagaConsumer {
	return &SagaConsumer{
		KafkaConsumer: kafkaConsumer,
		sagaService:   sagaService,
	}
}

func (c *SagaConsumer) Init() {
	logger := util.GetLogger().WithGroup("SagaConsumer.Init")

	if err := c.RegisterHandler(c.Handle); err != nil {
		logger.Error(constant.ErrTxtRegisterHandler, "err", err)
		panic(constant.ErrTxtRegisterHandler)
	}
}

func (c *SagaConsumer) Handle(ctx context.Context, r *consumerport.Message) error {
	logger := util.GetLogger().WithGroup("SagaConsumer.Handle")

	msg := &message.Message{}
	if err := json.Unmarshal(r.Value, msg); err != nil {
		logger.Error(constant.ErrTxtUnMarshalCommand, "value", r.Value, "err", err)
	}

	logger.Info("Received command", "command", msg)
	ctx, span := o11y.BeginSpanWithTraceID(ctx, msg.CorrelationID, msg.ParentID, "SagaConsumer")
	defer span.End()
	span.SetAttributes(
		o11y.AttrString("msg", fmt.Sprintf("%v", msg)),
	)

	switch msg.Name {
	case "StartSaga":
		cmd := &command.StartSaga{}
		if err := json.Unmarshal(r.Value, cmd); err != nil {
			logger.Error(constant.ErrTxtUnMarshalCommand, "err", err)
			span.RecordError(err)
			span.SetStatus(o11y.GetStatus(err))
			return err
		}
		return c.sagaService.Start(ctx, cmd)
	case "CarBooked":
		evt := &event.CarBooked{}
		if err := json.Unmarshal(r.Value, evt); err != nil {
			logger.Error(constant.ErrTxtUnMarshalEvent, "err", err)
			span.RecordError(err)
			span.SetStatus(o11y.GetStatus(err))
			return err
		}
		return c.sagaService.ProcessCarBooking(ctx, evt)
	case "CarBookingCanceled":
		evt := &event.CarBookingCanceled{}
		if err := json.Unmarshal(r.Value, evt); err != nil {
			logger.Error(constant.ErrTxtUnMarshalEvent, "err", err)
			span.RecordError(err)
			span.SetStatus(o11y.GetStatus(err))
			return err
		}
		if err := c.sagaService.CompensateCarBooking(ctx, evt); err != nil {
			logger.Error("Failed to compensate car booking", "err", err)
			span.RecordError(err)
			span.SetStatus(o11y.GetStatus(err))
			return err
		}
		return c.sagaService.MarkAbort(ctx, evt.Body.TripID)
	case "HotelBooked":
		evt := &event.HotelBooked{}
		if err := json.Unmarshal(r.Value, evt); err != nil {
			logger.Error(constant.ErrTxtUnMarshalEvent, "err", err)
			span.RecordError(err)
			span.SetStatus(o11y.GetStatus(err))
			return err
		}
		return c.sagaService.ProcessHotelBooking(ctx, evt)
	case "HotelBookingCanceled":
		evt := &event.HotelBookingCanceled{}
		if err := json.Unmarshal(r.Value, evt); err != nil {
			logger.Error(constant.ErrTxtUnMarshalEvent, "err", err)
			span.RecordError(err)
			span.SetStatus(o11y.GetStatus(err))
			return err
		}
		return c.sagaService.CompensateHotelBooking(ctx, evt)
	case "FlightBooked":
		evt := &event.FlightBooked{}
		if err := json.Unmarshal(r.Value, evt); err != nil {
			logger.Error(constant.ErrTxtUnMarshalEvent, "err", err)
			span.RecordError(err)
			span.SetStatus(o11y.GetStatus(err))
			return err
		}
		return c.sagaService.ProcessFlightBooking(ctx, evt)
	case "FlightBookingCanceled":
		evt := &event.FlightBookingCanceled{}
		if err := json.Unmarshal(r.Value, evt); err != nil {
			logger.Error(constant.ErrTxtUnMarshalEvent, "err", err)
			span.RecordError(err)
			span.SetStatus(o11y.GetStatus(err))
			return err
		}
		return c.sagaService.CompensateFlightBooking(ctx, evt)
	case "EndSaga":
		cmd := &command.EndSaga{}
		if err := json.Unmarshal(r.Value, cmd); err != nil {
			logger.Error(constant.ErrTxtUnMarshalCommand, "err", err)
			span.RecordError(err)
			span.SetStatus(o11y.GetStatus(err))
			return err
		}
		return c.sagaService.End(ctx, cmd)
	case "AbortSaga":
		cmd := &command.AbortSaga{}
		if err := json.Unmarshal(r.Value, cmd); err != nil {
			logger.Error(constant.ErrTxtUnMarshalCommand, "err", err)
			span.RecordError(err)
			span.SetStatus(o11y.GetStatus(err))
			return err
		}
		return c.sagaService.Abort(ctx, cmd)
	default:
		logger.Error(constant.ErrTxtUnknownCommand, "message", msg)
		err := errors.New(constant.ErrTxtUnknownCommand)
		span.RecordError(err)
		span.SetStatus(o11y.GetStatus(err))
		return err
	}
}
