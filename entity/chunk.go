package entity

type DocumentChunk struct {
	ID         string
	DocumentID string
	ChunkIndex int
	Content    string
	StartPage  int
	EndPage    int
	CharCount  int
}
