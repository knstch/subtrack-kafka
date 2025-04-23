package consumer

import (
	"context"
	"fmt"
	"log"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v3/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/knstch/subtrack-kafka/topics"
)

type Consumer struct {
	sub    *kafka.Subscriber
	router *message.Router
}

func NewConsumer(brokerAddr, group string, lg watermill.LoggerAdapter) (*Consumer, error) {
	sub, err := kafka.NewSubscriber(
		kafka.SubscriberConfig{
			Brokers:       []string{brokerAddr},
			Unmarshaler:   kafka.DefaultMarshaler{},
			ConsumerGroup: group,
		},
		lg,
	)
	if err != nil {
		return nil, err
	}

	router, err := message.NewRouter(message.RouterConfig{}, lg)
	if err != nil {
		log.Fatalf("router error: %v", err)
	}

	return &Consumer{
		sub:    sub,
		router: router,
	}, nil
}

func (c *Consumer) AddHandler(topic topics.KafkaTopic, handle message.NoPublishHandlerFunc) {
	c.router.AddNoPublisherHandler(fmt.Sprintf("handle-%s", topic), topic.String(), c.sub, handle)
}

func (c *Consumer) Run(ctx context.Context) error {
	if err := c.router.Run(ctx); err != nil {
		return err
	}

	return nil
}
