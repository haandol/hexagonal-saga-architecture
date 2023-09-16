package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/haandol/hexagonal/pkg/adapter/secondary/repository"
	"github.com/haandol/hexagonal/pkg/connector/producer"
	"github.com/haandol/hexagonal/pkg/dto"
	"github.com/haandol/hexagonal/pkg/message"
	"github.com/haandol/hexagonal/pkg/port/secondaryport/producerport"
	"github.com/haandol/hexagonal/pkg/port/secondaryport/repositoryport"
	"github.com/haandol/hexagonal/pkg/util"
	"github.com/haandol/hexagonal/pkg/util/o11y"
)

type MessageRelayService struct {
	kafkaProducer    producerport.Producer
	outboxRepository repositoryport.OutboxRepository
}

func NewMessageRelayService(
	kafkaProducer *producer.KafkaProducer,
	outboxRepository *repository.OutboxRepository,
) *MessageRelayService {
	return &MessageRelayService{
		kafkaProducer:    kafkaProducer,
		outboxRepository: outboxRepository,
	}
}

func (s *MessageRelayService) Fetch(ctx context.Context, batchSize int) ([]dto.Outbox, error) {
	logger := util.GetLogger().WithGroup("MessageRelayService.Fetch")

	// TODO: group by kafkaKey and send them parallell
	messages, err := s.outboxRepository.QueryUnsent(ctx, batchSize)
	if err != nil {
		logger.Error("failed to query unsent messages", "err", err)
		return nil, err
	}

	return messages, nil
}

func (s *MessageRelayService) Relay(ctx context.Context, messages []dto.Outbox) error {
	logger := util.GetLogger().WithGroup("MessageRelayService.Relay")

	var wg sync.WaitGroup
	sentIDs := make([]uint, len(messages))

	for _, msg := range messages {
		wg.Add(1)
		go func(m dto.Outbox) {
			defer wg.Done()

			var msg message.Message
			data := []byte(m.KafkaValue)
			if err := json.Unmarshal(data, &msg); err != nil {
				logger.Error("failed to unmarshal message", "err", err)
				return
			}

			ctx, span := o11y.BeginSpanWithTraceID(ctx, msg.CorrelationID, msg.ParentID, "Relay")
			defer span.End()

			span.SetAttributes(
				o11y.AttrInt("msgId", int(m.ID)),
				o11y.AttrString("msg", fmt.Sprintf("%v", msg)),
			)

			if err := s.kafkaProducer.Produce(ctx, m.KafkaTopic, m.KafkaKey, []byte(m.KafkaValue)); err != nil {
				logger.Error("failed to produce message", "err", err)
				span.RecordError(err)
				span.SetStatus(o11y.GetStatus(err))
				return
			}
			sentIDs = append(sentIDs, m.ID)
		}(msg)
	}
	wg.Wait()

	if len(messages) > 0 {
		logger.Info("sent messages", "total", len(messages), "sent", len(sentIDs))
	}

	if len(sentIDs) > 0 {
		if err := s.outboxRepository.MarkSentInBatch(ctx, sentIDs); err != nil {
			logger.Error("failed to mark message as sent", "err", err)
			return err
		}
	}

	return nil
}
