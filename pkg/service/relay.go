package service

import (
	"context"
	"sync"

	"github.com/haandol/hexagonal/pkg/adapter/secondary/repository"
	"github.com/haandol/hexagonal/pkg/connector/producer"
	"github.com/haandol/hexagonal/pkg/dto"
	"github.com/haandol/hexagonal/pkg/port/secondaryport/producerport"
	"github.com/haandol/hexagonal/pkg/port/secondaryport/repositoryport"
	"github.com/haandol/hexagonal/pkg/util"
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
	logger := util.GetLogger().With(
		"module", "MessageRelayService",
		"func", "Fetch",
	)

	// TODO: group by kafkaKey and send them parallell
	messages, err := s.outboxRepository.QueryUnsent(ctx, batchSize)
	if err != nil {
		logger.Errorw("failed to query unsent messages", "err", err.Error())
		return nil, err
	}

	return messages, nil
}

func (s *MessageRelayService) Relay(ctx context.Context, messages []dto.Outbox) error {
	logger := util.GetLogger().WithContext(ctx).With(
		"module", "MessageRelayService",
		"func", "Relay",
	)

	var sentIDs []uint
	var wg sync.WaitGroup

	for _, msg := range messages {
		wg.Add(1)
		go func(m dto.Outbox) {
			defer wg.Done()
			if err := s.kafkaProducer.Produce(ctx, m.KafkaTopic, m.KafkaKey, []byte(m.KafkaValue)); err != nil {
				logger.Errorw("failed to produce message", "err", err)
				return
			}
			sentIDs = append(sentIDs, m.ID)
		}(msg)
	}
	wg.Wait()

	if len(messages) > 0 {
		logger.Infow("sent messages", "total", len(messages), "sent", len(sentIDs))
	}

	if len(sentIDs) > 0 {
		if err := s.outboxRepository.MarkSentInBatch(ctx, sentIDs); err != nil {
			logger.Errorw("failed to mark message as sent", "err", err)
			return err
		}
	}

	return nil
}
