package config

import (
	"context"

	"gitlab.spesolution.net/bni-merchant-management-system/go-sekeleton/internal/queue"
)

func NewRabbitMQInstance(ctx context.Context, cfg *RabbitMQOption) (*queue.RabbitMQ, error) {
	rabbit := &queue.RabbitMQ{
		Ctx:        ctx,
		Uri:        cfg.Uri,
		Exchange:   cfg.Exchange,
		Kind:       cfg.QueueType,
		Prefix:     cfg.QueuePrefix,
		RetryCount: cfg.QueueRetryCount,
		Err:        make(chan error),
	}

	if err := rabbit.Connect(); err != nil {
		return nil, err
	}

	return rabbit, nil
}
