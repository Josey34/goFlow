package entity

type ProcessingStats struct {
	TotalProcessed    int
	DuplicatesFound   int
	ErrorsEncountered int
	AvgProcessingTime float64
}
