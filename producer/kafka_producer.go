package producer

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"

	kafkaPkg "github.com/knstch/subtrack-kafka/topics"
)

type Producer struct {
	Addr     string
	Balancer *kafka.LeastBytes
}

func NewProducer(addr string) *Producer {
	return &Producer{
		Addr:     addr,
		Balancer: &kafka.LeastBytes{},
	}
}

func (p *Producer) SendMessage(topic kafkaPkg.KafkaTopic, key string, value interface{}) error {
	body, err := json.Marshal(value)
	if err != nil {
		return err
	}

	writer := kafka.Writer{
		Addr:     kafka.TCP(p.Addr),
		Topic:    topic.String(),
		Balancer: p.Balancer,
	}

	defer writer.Close()

	err = writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(key),
			Value: body,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
