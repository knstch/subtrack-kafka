package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ThreeDotsLabs/watermill/message"
)

type TypedHandler[T any] func(ctx context.Context, payload T) error

func JSONHandler[T any](handler TypedHandler[T]) func(msg *message.Message) error {
	return func(msg *message.Message) error {
		var payload T
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			return fmt.Errorf("failed to unmarshal message: %w", err)
		}
		return handler(context.Background(), payload)
	}
}
