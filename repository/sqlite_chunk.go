package repository

import (
	"context"
	"database/sql"

	"goflow/entity"
)

type SQLiteChunkRepo struct {
	db *sql.DB
}

func NewSQLiteChunkRepo(db *sql.DB) *SQLiteChunkRepo {
	return &SQLiteChunkRepo{db: db}
}

func (r *SQLiteChunkRepo) InsertBatch(ctx context.Context, chunks []entity.DocumentChunk) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, chunk := range chunks {
		query := `
			INSERT INTO document_chunks 
			(id, document_id, chunk_index, content, start_page, end_page, char_count)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`
		_, err := tx.ExecContext(ctx, query,
			chunk.ID, chunk.DocumentID, chunk.ChunkIndex, chunk.Content,
			chunk.StartPage, chunk.EndPage, chunk.CharCount,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *SQLiteChunkRepo) FindByDocID(ctx context.Context, docID string) ([]entity.DocumentChunk, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id, document_id, chunk_index, content, start_page, end_page, char_count FROM document_chunks WHERE document_id = ? ORDER BY chunk_index",
		docID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chunks []entity.DocumentChunk
	for rows.Next() {
		var chunk entity.DocumentChunk
		err := rows.Scan(&chunk.ID, &chunk.DocumentID, &chunk.ChunkIndex, &chunk.Content, &chunk.StartPage, &chunk.EndPage, &chunk.CharCount)
		if err != nil {
			return nil, err
		}
		chunks = append(chunks, chunk)
	}
	return chunks, rows.Err()
}

func (r *SQLiteChunkRepo) Search(ctx context.Context, query string) ([]entity.DocumentChunk, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id, document_id, chunk_index, content, start_page, end_page, char_count FROM document_chunks WHERE content LIKE ?",
		"%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chunks []entity.DocumentChunk
	for rows.Next() {
		var chunk entity.DocumentChunk
		err := rows.Scan(&chunk.ID, &chunk.DocumentID, &chunk.ChunkIndex, &chunk.Content, &chunk.StartPage, &chunk.EndPage, &chunk.CharCount)
		if err != nil {
			return nil, err
		}
		chunks = append(chunks, chunk)
	}
	return chunks, rows.Err()
}
