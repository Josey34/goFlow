package worker

import (
	"goflow/entity"
	"time"
)

type ProcessingTask struct {
	Event         *entity.Event
	CreatedAt     time.Time
	RetryCount    int
	CorrelationID string
}

func NewProcessingTask(event *entity.Event) *ProcessingTask {
	return &ProcessingTask{
		Event:         event,
		CreatedAt:     time.Now(),
		RetryCount:    0,
		CorrelationID: event.DocumentID,
	}
}
