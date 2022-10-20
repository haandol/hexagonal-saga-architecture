package producer

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/haandol/hexagonal/pkg/config"
	"github.com/haandol/hexagonal/pkg/util"
	"github.com/twmb/franz-go/pkg/kgo"
)

type KafkaProducer struct {
	client *kgo.Client
	pool   *sync.Pool
}

func NewKafkaProducer(cfg *config.Config) *KafkaProducer {
	opts := BuildProducerOpts(cfg.Kafka.Seeds)
	if strings.Contains(cfg.Kafka.Seeds[0], "9094") {
		opts = append(opts, kgo.DialTLSConfig(new(tls.Config)))
	}

	client, err := kgo.NewClient(opts...)
	if err != nil {
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

func BuildProducerOpts(seeds []string) []kgo.Opt {
	return []kgo.Opt{
		kgo.SeedBrokers(seeds...),
		kgo.ProducerLinger(0),
		kgo.AllowAutoTopicCreation(), // for dev only
		kgo.ProducerBatchCompression(kgo.GzipCompression()),
		kgo.WithLogger(kgo.BasicLogger(os.Stderr, kgo.LogLevelInfo, func() string {
			return fmt.Sprintf("%s\t", time.Now().Format(time.RFC3339))
		})),
	}
}

func (p KafkaProducer) newRecord(topic, key string, val []byte) *kgo.Record {
	r := p.pool.Get().(*kgo.Record)
	r.Topic = topic
	r.Key = []byte(key)
	r.Value = val
	return r
}

func (p *KafkaProducer) Produce(ctx context.Context, topic, key string, val []byte) error {
	logger := util.GetLogger().WithContext(ctx).With(
		"module", "KafkaProducer",
		"func", "Produce",
	)
	logger.Info("Producing...")

	r := p.newRecord(topic, key, val)
	if err := p.client.ProduceSync(ctx, r).FirstErr(); err != nil {
		logger.Errorw("produce failed", "err", err.Error())
		return err
	}
	logger.Debugw("Message produced", "record", r)
	p.pool.Put(r)
	return nil
}

func (p *KafkaProducer) Close(ctx context.Context) error {
	logger := util.GetLogger().With(
		"module", "KafkaProducer",
		"func", "Close",
	)
	logger.Info("Closing...")

	if err := p.client.Flush(ctx); err != nil {
		logger.Errorw("failed to flush", "err", err.Error())
	}

	p.client.Close()
	return nil
}
