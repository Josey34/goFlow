package mock

import (
	"context"
	"errors"
	"goflow/entity"
)

type MockConsumer struct {
	Events      []*entity.Event
	callCount   int
	shouldError bool
}

func NewMockConsumer() *MockConsumer {
	return &MockConsumer{
		Events: make([]*entity.Event, 0),
	}
}

func (mc *MockConsumer) Consume(ctx context.Context) (*entity.Event, error) {
	if mc.shouldError && mc.callCount == 0 {
		mc.callCount++
		return nil, errors.New("consume failed")
	}

	if mc.callCount >= len(mc.Events) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	event := mc.Events[mc.callCount]
	mc.callCount++
	return event, nil
}

func (mc *MockConsumer) Acknowledge(ctx context.Context, messageID string) error {
	return nil
}

func (mc *MockConsumer) WithError() *MockConsumer {
	mc.shouldError = true
	return mc
}

func (mc *MockConsumer) WithEvents(events ...*entity.Event) *MockConsumer {
	mc.Events = append(mc.Events, events...)
	return mc
}

func (mc *MockConsumer) Close() error {
	return nil
}
