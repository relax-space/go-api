package adapters

import (
	"context"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/pangpanglabs/goutils/behaviorlog"
	"github.com/pangpanglabs/goutils/kafka"
	"github.com/sirupsen/logrus"
)

var MessagePublisher *MessagePublish
var messagePublisherOnce sync.Once

const (
	EventCreatedFruit = "FruitCreated"
)

type MessagePublish struct {
	producer *kafka.Producer
}

func NewMessagePublisher(kafkaConfig kafka.Config) *MessagePublish {
	messagePublisherOnce.Do(func() {
		producer, err := kafka.NewProducer(kafkaConfig.Brokers, kafkaConfig.Topic, func(c *sarama.Config) {
			c.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
			c.Producer.Compression = sarama.CompressionGZIP     // Compress messages
			c.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms
		})

		if err != nil {
			logrus.Error("Fail to producer create")
			return
		}

		MessagePublisher = &MessagePublish{
			producer: producer,
		}
	})
	return MessagePublisher
}

func (publisher *MessagePublish) Close() {
	publisher.producer.Close()
}

func (publisher *MessagePublish) Publish(ctx context.Context, payload interface{}, status string) error {
	m := map[string]interface{}{
		"authToken": behaviorlog.FromCtx(ctx).AuthToken,
		"requestId": behaviorlog.FromCtx(ctx).RequestID,
		"status":    status,
		"payload":   payload,
		"createdAt": time.Now().UTC(),
	}
	return publisher.producer.Send(m)
}
