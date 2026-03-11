package entity

import "time"

type ProcessingResult struct {
	ID            string
	DocumentID    string
	ExtractedText string
	PageCount     int
	FileHash      string
	IsDuplicate   bool
	ThumbnailInfo string
	ProcessedAt   time.Time
	ErrorMessage  string `json:"error_message,omitempty"`
}
