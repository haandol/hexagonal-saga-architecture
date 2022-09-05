package publisher

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/haandol/hexagonal/pkg/message"
	"github.com/haandol/hexagonal/pkg/message/command"
	"github.com/haandol/hexagonal/pkg/message/event"
	"github.com/haandol/hexagonal/pkg/port/secondaryport/producerport"
	"github.com/haandol/hexagonal/pkg/util"
)

func PublishSagaEnded(ctx context.Context, p producerport.Producer, cmd *command.EndSaga) error {
	logger := util.GetLogger().With(
		"module", "Publisher",
		"func", "PublishSagaEnded",
	)

	evt := &event.SagaEnded{
		Message: message.Message{
			Name:          "SagaEnded",
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: cmd.CorrelationID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: event.SagaEndedBody{
			SagaID: cmd.Body.SagaID,
		},
	}
	if err := util.ValidateStruct(evt); err != nil {
		return err
	}
	v, err := json.Marshal(evt)
	if err != nil {
		logger.Errorw("failed to marshal saga aborted event", "event", evt, "err", err.Error())
	}

	if err := p.Produce(ctx, "trip-service", cmd.CorrelationID, v); err != nil {
		logger.Errorw("failed to produce saga ended event", "event", evt, "err", err.Error())
		return err
	}

	return nil
}

func PublishSagaAborted(ctx context.Context, p producerport.Producer, cmd *command.AbortSaga) error {
	logger := util.GetLogger().With(
		"module", "Publisher",
		"func", "PublishSagaEnded",
	)

	evt := &event.SagaAborted{
		Message: message.Message{
			Name:          "SagaAborted",
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: cmd.CorrelationID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: event.SagaAbortedBody{
			SagaID: cmd.Body.SagaID,
			Reason: cmd.Body.Reason,
			Source: cmd.Body.Source,
		},
	}
	if err := util.ValidateStruct(evt); err != nil {
		return err
	}
	v, err := json.Marshal(evt)
	if err != nil {
		logger.Errorw("failed to marshal saga aborted event", "event", evt, "err", err.Error())
	}

	if err := p.Produce(ctx, "trip-service", cmd.CorrelationID, v); err != nil {
		logger.Errorw("failed to produce saga aborted event", "event", evt, "err", err.Error())
		return err
	}

	return nil
}
