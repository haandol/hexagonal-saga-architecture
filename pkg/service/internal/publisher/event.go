package publisher

import (
	"context"
	"encoding/json"

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

	event := &event.SagaAborted{
		Message: message.Message{
			Name: "SagaEnded",
		},
		Body: event.SagaAbortedBody{
			SagaID: cmd.Body.SagaID,
		},
	}
	v, err := json.Marshal(event)
	if err != nil {
		logger.Errorw("failed to marshal saga aborted event", "event", event, "err", err.Error())
	}

	if err := p.Produce(ctx, "trip-service", cmd.CorrelationID, v); err != nil {
		logger.Errorw("failed to produce saga ended event", "event", event, "err", err.Error())
		return err
	}

	return nil
}

func PublishSagaAborted(ctx context.Context, p producerport.Producer, cmd *command.AbortSaga) error {
	logger := util.GetLogger().With(
		"module", "Publisher",
		"func", "PublishSagaEnded",
	)

	event := &event.SagaAborted{
		Message: message.Message{
			Name: "SagaEnded",
		},
		Body: event.SagaAbortedBody{
			SagaID: cmd.Body.SagaID,
		},
	}
	v, err := json.Marshal(event)
	if err != nil {
		logger.Errorw("failed to marshal saga aborted event", "event", event, "err", err.Error())
	}

	if err := p.Produce(ctx, "trip-service", cmd.CorrelationID, v); err != nil {
		logger.Errorw("failed to produce saga aborted event", "event", event, "err", err.Error())
		return err
	}

	return nil
}
