package consumer

import (
	"context"
	"crypto/tls"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/haandol/hexagonal/pkg/config"
	"github.com/haandol/hexagonal/pkg/port/primaryport/consumerport"
	"github.com/haandol/hexagonal/pkg/util"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/plugin/kzap"
)

type KafkaConsumer struct {
	client        *kgo.Client
	topic         string
	messageExpiry time.Duration
	handler       consumerport.HandlerFunc
	batchSize     int
}

const ConsumerTimeout = 30 * time.Second

func NewKafkaConsumer(cfg *config.Kafka, groupID, topic string) *KafkaConsumer {
	opts := buildConsumerOpts(cfg.Seeds, groupID, topic)
	if strings.Contains(cfg.Seeds[0], "9094") {
		opts = append(opts, kgo.DialTLSConfig(new(tls.Config)))
	}

	client, err := kgo.NewClient(opts...)
	if err != nil {
		log.Panic(err.Error())
	}

	return &KafkaConsumer{
		client:        client,
		topic:         topic,
		messageExpiry: time.Duration(cfg.MessageExpirySec) * time.Second,
		batchSize:     cfg.BatchSize,
		handler:       nil,
	}
}

func buildConsumerOpts(seeds []string, group, topic string) []kgo.Opt {
	return []kgo.Opt{
		kgo.SeedBrokers(seeds...),
		kgo.ConsumerGroup(group),
		kgo.ConsumeTopics(topic),
		kgo.DisableAutoCommit(),
		kgo.AllowAutoTopicCreation(), // TODO: only for the dev
		kgo.WithLogger(kzap.New(
			util.GetLogger().With("package", "consumer").Desugar(),
			kzap.Level(kgo.LogLevelInfo),
		)),
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
		logger.Panic("handler not registered")
	}

	for {
		logger.Infow("Polling...", "topic", c.topic)
		ctx := context.Background()

		fetches := c.client.PollRecords(ctx, c.batchSize)
		if fetches.IsClientClosed() {
			logger.Infow("Client closed", "topic", c.topic)
			return
		}
		if errs := fetches.Errors(); len(errs) > 0 {
			logger.Panicw("failed to fetch records", "topic", c.topic, "err", errs[0].Err.Error())
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
			if c.messageExpiry > 0 && time.Since(record.Timestamp) > c.messageExpiry {
				logger.Warnw("message expired", "expirySec", c.messageExpiry, "key", key)
				return
			}

			if err := c.handler(ctx, message); err != nil {
				logger.Panicw("error handling message", "err", err.Error())
			}
		})

		if err := c.client.CommitUncommittedOffsets(ctx); err != nil {
			logger.Panicw("failed to commit offsets", "topic", c.topic, "err", err.Error())
		}
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
