package consumer

import (
	"context"
	"encoding/json"

	"github.com/haandol/hexagonal/message/command"
	"github.com/haandol/hexagonal/pkg/port/primaryport/consumerport"
	"github.com/haandol/hexagonal/pkg/util"
)

type SagaConsumer struct {
	*KafkaConsumer
}

func NewSagaConsumer(kafkaConsumer *KafkaConsumer) *SagaConsumer {
	return &SagaConsumer{
		KafkaConsumer: kafkaConsumer,
	}
}

func (c *SagaConsumer) Init() {
	logger := util.GetLogger().With(
		"module", "SagaConsumer",
		"func", "Init",
	)

	if err := c.RegisterHandler(c.Handle); err != nil {
		logger.Fatalw("Failed to register handler", "err", err.Error())
	}
}

func (c *SagaConsumer) Handle(ctx context.Context, r *consumerport.Message) error {
	logger := util.GetLogger().With(
		"module", "SagaConsumer",
		"func", "Handle",
	)

	cmd := &command.Command{}
	if err := json.Unmarshal(r.Value, cmd); err != nil {
		logger.Errorw("Failed to unmarshal command", "err", err.Error())
	}

	logger.Infow("Received command", "command", cmd)

	return nil
}
