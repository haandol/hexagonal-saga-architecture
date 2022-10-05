package service

import (
	"context"

	"github.com/haandol/hexagonal/pkg/port/secondaryport/producerport"
	"github.com/haandol/hexagonal/pkg/port/secondaryport/repositoryport"
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

func (s *MessageRelayService) Relay(ctx context.Context) (int, error) {
	var num int

	// TODO: group by kafkaKey and send them parallell
	messages, err := s.outboxRepository.QueryUnsent(ctx)
	if err != nil {
		return 0, err
	}

	for i, m := range messages {
		if err := s.producer.Produce(ctx, m.KafkaTopic, m.KafkaKey, []byte(m.KafkaValue)); err != nil {
			return num, err
		}

		if err := s.outboxRepository.Delete(ctx, m.ID); err != nil {
			return num, err
		}

		num = i
	}

	return num, nil
}
