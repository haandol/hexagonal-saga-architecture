package producer

import (
	"context"
	"crypto/tls"
	"strings"
	"sync"

	"github.com/haandol/hexagonal/pkg/config"
	"github.com/haandol/hexagonal/pkg/util"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/plugin/kzap"
)

var (
	kafkaProducer *KafkaProducer
)

type KafkaProducer struct {
	client *kgo.Client
	pool   *sync.Pool
}

func Connect(cfg *config.Kafka) (*KafkaProducer, error) {
	if kafkaProducer == nil {
		kafkaProducer = newKafkaProducer(cfg)
	}

	if err := kafkaProducer.client.Ping(context.Background()); err != nil {
		return nil, err
	}

	return kafkaProducer, nil
}

func newKafkaProducer(cfg *config.Kafka) *KafkaProducer {
	opts := buildProducerOpts(cfg.Seeds)
	if strings.Contains(cfg.Seeds[0], "9094") {
		opts = append(opts, kgo.DialTLSConfig(new(tls.Config)))
	}

	client, err := kgo.NewClient(opts...)
	if err != nil {
		panic(err)
	}

	if err := client.Ping(context.Background()); err != nil {
		panic(err)
	}

	recordBytes := 4096
	pool := sync.Pool{New: func() any {
		return kgo.SliceRecord(make([]byte, recordBytes))
	}}

	return &KafkaProducer{
		client: client,
		pool:   &pool,
	}
}

func buildProducerOpts(seeds []string) []kgo.Opt {
	return []kgo.Opt{
		kgo.SeedBrokers(seeds...),
		kgo.AllowAutoTopicCreation(), // for dev only
		kgo.RecordPartitioner(kgo.StickyKeyPartitioner(nil)),
		kgo.ProducerBatchCompression(kgo.GzipCompression()),
		kgo.WithLogger(kzap.New(
			util.GetLogger().With("package", "producer").Desugar(),
			kzap.Level(kgo.LogLevelInfo),
		)),
	}
}

func (p *KafkaProducer) newRecord(topic, key string, val []byte) *kgo.Record {
	r := p.pool.Get().(*kgo.Record)
	r.Topic = topic
	r.Key = []byte(key)
	r.Value = val
	return r
}

func (p *KafkaProducer) Produce(ctx context.Context, topic, key string, val []byte) error {
	logger := util.GetLogger().With(
		"module", "KafkaProducer",
		"func", "Produce",
	)
	logger.Info("Producing...")

	r := p.newRecord(topic, key, val)
	if err := p.client.ProduceSync(ctx, r).FirstErr(); err != nil {
		logger.Errorw("produce failed", "err", err)
		return err
	}
	logger.Debugw("Message produced", "sent bytes", len(r.Value))
	p.pool.Put(r)
	return nil
}

func Close(ctx context.Context) error {
	logger := util.GetLogger().With(
		"module", "KafkaProducer",
		"func", "Close",
	)
	logger.Info("Closing producer...")

	if kafkaProducer == nil {
		logger.Warn("Producer is not connected")
		return nil
	}

	if err := kafkaProducer.client.Flush(ctx); err != nil {
		logger.Errorw("failed to flush", "err", err)
		return err
	}

	kafkaProducer.client.Close()

	return nil
}
