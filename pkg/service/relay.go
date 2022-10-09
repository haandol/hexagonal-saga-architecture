package service

import (
	"context"
	"sync"

	"github.com/haandol/hexagonal/pkg/dto"
	"github.com/haandol/hexagonal/pkg/port/secondaryport/producerport"
	"github.com/haandol/hexagonal/pkg/port/secondaryport/repositoryport"
	"github.com/haandol/hexagonal/pkg/util"
)

type MessageRelayService struct {
	producer         producerport.Producer
	outboxRepository repositoryport.OutboxRepository
}

func NewMessageRelayService(
	producer producerport.Producer,
	outboxRepository repositoryport.OutboxRepository,
) *MessageRelayService {
	return &MessageRelayService{
		producer:         producer,
		outboxRepository: outboxRepository,
	}
}

func (s *MessageRelayService) Relay(ctx context.Context, batchSize int) error {
	logger := util.GetLogger().With(
		"module", "MessageRelayService",
		"func", "Relay",
	)

	var numSent int

	// TODO: group by kafkaKey and send them parallell
	messages, err := s.outboxRepository.QueryUnsent(ctx, batchSize)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	for _, msg := range messages {
		wg.Add(1)
		go func(m dto.Outbox) {
			defer wg.Done()
			if err := s.producer.Produce(ctx, m.KafkaTopic, m.KafkaKey, []byte(m.KafkaValue)); err != nil {
				logger.Errorw("failed to produce message", "err", err)
				return
			}

			if err := s.outboxRepository.MarkSent(ctx, m.ID); err != nil {
				logger.Errorw("failed to mark message as sent", "err", err)
				return
			}
			numSent++
		}(msg)
	}
	wg.Wait()

	if numSent > 0 {
		logger.Infow("sent messages", "numSent", numSent)
	}

	return nil
}
