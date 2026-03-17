package service

import (
	"context"
	"goflow/entity"
)

type EventConsumer interface {
	Consume(ctx context.Context) (*entity.Event, error)
	Acknowledge(ctx context.Context, messageID string) error
	Close() error
}
