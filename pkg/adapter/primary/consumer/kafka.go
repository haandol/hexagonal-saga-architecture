package consumer

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/haandol/hexagonal/pkg/config"
	"github.com/haandol/hexagonal/pkg/port/primaryport/consumerport"
	"github.com/haandol/hexagonal/pkg/util"
	"github.com/twmb/franz-go/pkg/kgo"
)

type KafkaConsumer struct {
	client           *kgo.Client
	topic            string
	messageExpirySec time.Duration
	handler          consumerport.HandlerFunc
	batchSize        int
}

func NewKafkaConsumer(cfg config.Kafka, groupID string, topic string) *KafkaConsumer {
	opts := BuildConsumerOpts(cfg.Seeds, groupID, topic)
	client, err := kgo.NewClient(opts...)
	if err != nil {
		log.Fatal(err.Error())
	}

	return &KafkaConsumer{
		client:           client,
		topic:            topic,
		messageExpirySec: time.Duration(cfg.MessageExpirySec) * time.Second,
		batchSize:        cfg.BatchSize,
		handler:          nil,
	}
}

func BuildConsumerOpts(seeds []string, group, topic string) []kgo.Opt {
	return []kgo.Opt{
		kgo.SeedBrokers(seeds...),
		kgo.ConsumerGroup(group),
		kgo.ConsumeTopics(topic),
		kgo.WithLogger(kgo.BasicLogger(os.Stderr, kgo.LogLevelInfo, func() string {
			return fmt.Sprintf("%s\t", time.Now().Format(time.RFC3339))
		})),
	}
}

func (c *KafkaConsumer) RegisterHandler(h consumerport.HandlerFunc) error {
	logger := util.GetLogger().With(
		"module", "KafkaConsumer",
		"func", "RegisterHandler",
	)
	logger.Info("Registering handler...")

	if c.handler != nil {
		logger.Error("handler already registered")
		return errors.New("handler already registered")
	}

	c.handler = h
	logger.Info("registered handler")

	return nil
}

// Consume - consume messages from Kafka and dispatch to handlers
func (c *KafkaConsumer) Consume() {
	logger := util.GetLogger().With(
		"module", "KafkaConsumer",
		"func", "Consume",
	)
	logger.Infow("Consuming Topic", "topic", c.topic)

	// check initialized
	if c.handler == nil {
		logger.Fatal("handler not registered")
	}

	for {
		logger.Infow("Polling...", "topic", c.topic)
		fetches := c.client.PollRecords(context.Background(), c.batchSize)
		if fetches.IsClientClosed() {
			return
		}
		if errs := fetches.Errors(); len(errs) > 0 {
			logger.Fatal(errs[0].Err.Error())
		}

		fetches.EachRecord(func(record *kgo.Record) {
			key := string(record.Key)
			logger.Infow("Message received", "key", key)

			message := &consumerport.Message{
				Topic:     record.Topic,
				Key:       key,
				Value:     record.Value,
				Timestamp: record.Timestamp,
			}
			if c.messageExpirySec > 0 && time.Since(record.Timestamp) > c.messageExpirySec {
				logger.Infow("Message expired", "expirySec", c.messageExpirySec, "key", key)
				return
			}

			err := c.handler(context.TODO(), message)
			if err != nil {
				logger.Errorw("Error handling message", "err", err.Error())
			}
		})
	}
}

func (c *KafkaConsumer) Close(ctx context.Context) error {
	logger := util.GetLogger().With(
		"module", "KafkaConsumer",
		"func", "Close",
	)
	logger.Infow("Closing...", "topic", c.topic)

	c.client.Close()
	return nil
}
