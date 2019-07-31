package adapters

import (
	"go-api/models"
	"nomni/utils/eventconsume"

	"github.com/pangpanglabs/goutils/kafka"
)

func Consume(serviceName string, kafkaConfig kafka.Config,
	f func(eventconsume.ConsumeContext) error, filters ...eventconsume.Filter) error {
	return eventconsume.NewEventConsumer(
		serviceName,
		kafkaConfig.Brokers,
		kafkaConfig.Topic,
		filters).Handle(f)
}

func EventFruit(c eventconsume.ConsumeContext) error {
	var fruit models.Fruit
	if err := c.Bind(&fruit); err != nil {
		return err
	}
	if c.Status() == "FruitCreated" {
		if _, err := (&fruit).Create(c.Context()); err != nil {
			return err
		}
	}
	return nil
}
