package entity

import "time"

type Event struct {
	DocumentID string    `json:"document_id"`
	Filename   string    `json:"filename"`
	UserID     string    `json:"user_id"`
	EventType  string    `json:"event_type"`
	Timestamp  time.Time `json:"timestamp"`
	MessageID  string    `json:"message_id,omitempty"`
}
