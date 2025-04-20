package producer

import (
	"context"
	"encoding/json"
	"net"

	"github.com/segmentio/kafka-go"

	kafkaPkg "github.com/knstch/subtrack-kafka/topics"
)

type Producer struct {
	addr     net.Addr
	balancer *kafka.LeastBytes
}

func NewProducer(addr string) *Producer {
	return &Producer{
		addr:     kafka.TCP(addr),
		balancer: &kafka.LeastBytes{},
	}
}

func (p *Producer) SendMessage(topic kafkaPkg.KafkaTopic, key string, value interface{}) error {
	body, err := json.Marshal(value)
	if err != nil {
		return err
	}

	writer := kafka.Writer{
		Addr:     p.addr,
		Topic:    topic.String(),
		Balancer: p.balancer,
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
