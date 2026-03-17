package repository

import (
	"context"
	"database/sql"
	"goflow/entity"
	"time"
)

type SQLiteResultRepo struct {
	db *sql.DB
}

func NewSQLiteResultRepo(db *sql.DB) *SQLiteResultRepo {
	return &SQLiteResultRepo{db: db}
}

func (r *SQLiteResultRepo) Insert(ctx context.Context, result *entity.ProcessingResult) error {
	query := `
		INSERT INTO processing_results 
		(id, document_id, extracted_text, page_count, file_hash, is_duplicate, thumbnail_info, processed_at, error_message)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		result.ID, result.DocumentID, result.ExtractedText, result.PageCount,
		result.FileHash, result.IsDuplicate, result.ThumbnailInfo,
		time.Now(), result.ErrorMessage,
	)
	return err
}

func (r *SQLiteResultRepo) FindByDocID(ctx context.Context, docID string) (*entity.ProcessingResult, error) {
	result := &entity.ProcessingResult{}
	err := r.db.QueryRowContext(ctx,
		"SELECT id, document_id, extracted_text, page_count, file_hash, is_duplicate, thumbnail_info, processed_at FROM processing_results WHERE document_id = ?",
		docID).
		Scan(&result.ID, &result.DocumentID, &result.ExtractedText, &result.PageCount, &result.FileHash, &result.IsDuplicate, &result.ThumbnailInfo, &result.ProcessedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return result, err
}

func (r *SQLiteResultRepo) FindByHash(ctx context.Context, hash string) (*entity.ProcessingResult, error) {
	result := &entity.ProcessingResult{}
	err := r.db.QueryRowContext(ctx,
		"SELECT id, document_id, extracted_text, page_count, file_hash, is_duplicate, thumbnail_info, processed_at FROM processing_results WHERE file_hash = ?",
		hash).
		Scan(&result.ID, &result.DocumentID, &result.ExtractedText, &result.PageCount, &result.FileHash, &result.IsDuplicate, &result.ThumbnailInfo, &result.ProcessedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return result, err
}
