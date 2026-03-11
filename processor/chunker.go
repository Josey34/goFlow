package processor

import (
	"strings"

	"goflow/entity"

	"github.com/google/uuid"
)

type ChunkConfig struct {
	ChunkSize    int
	ChunkOverlap int
}

func ChunkText(text string, docID string, pageCount int, cfg ChunkConfig) []entity.DocumentChunk {
	if text == "" || cfg.ChunkSize == 0 {
		return []entity.DocumentChunk{}
	}

	pageMap := buildPageMap(text, pageCount)

	chunks := []entity.DocumentChunk{}
	chunkIndex := 0
	step := cfg.ChunkSize - cfg.ChunkOverlap

	for i := 0; i < len(text); i += step {
		endIdx := i + cfg.ChunkSize
		if endIdx > len(text) {
			endIdx = len(text)
		}

		if endIdx-i < cfg.ChunkSize/2 && i > 0 {
			break
		}

		content := text[i:endIdx]
		startPage := pageMap[i]
		endPage := pageMap[endIdx-1]

		chunks = append(chunks, entity.DocumentChunk{
			ID:         uuid.New().String(),
			DocumentID: docID,
			ChunkIndex: chunkIndex,
			Content:    content,
			StartPage:  startPage,
			EndPage:    endPage,
			CharCount:  len(content),
		})

		chunkIndex++
		if endIdx >= len(text) {
			break
		}
	}

	return chunks
}

func buildPageMap(text string, pageCount int) map[int]int {
	pageMap := make(map[int]int)
	pages := strings.Split(text, "\n")
	charPos := 0
	pageNum := 1

	for _, pageText := range pages {
		for i := 0; i < len(pageText); i++ {
			pageMap[charPos] = pageNum
			charPos++
		}
		charPos++
		if pageNum < pageCount {
			pageNum++
		}
	}

	return pageMap
}
