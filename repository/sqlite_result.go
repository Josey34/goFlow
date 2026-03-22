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
		(id, document_id, text_content, created_at)
		VALUES (?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		result.ID, result.DocumentID, result.ExtractedText, time.Now(),
	)
	return err
}

func (r *SQLiteResultRepo) FindByDocID(ctx context.Context, docID string) (*entity.ProcessingResult, error) {
	result := &entity.ProcessingResult{}
	err := r.db.QueryRowContext(ctx,
		"SELECT id, document_id, text_content, created_at FROM processing_results WHERE document_id = ?",
		docID).
		Scan(&result.ID, &result.DocumentID, &result.ExtractedText, &result.ProcessedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return result, err
}

func (r *SQLiteResultRepo) FindByHash(ctx context.Context, hash string) (*entity.ProcessingResult, error) {
	return nil, nil
}

func (r *SQLiteResultRepo) GetStats(ctx context.Context) (*entity.ProcessingStats, error) {
	query := `
		SELECT
			COUNT(*) as total_processed
		FROM processing_results
	`

	stats := &entity.ProcessingStats{}
	err := r.db.QueryRowContext(ctx, query).Scan(&stats.TotalProcessed)

	if err == sql.ErrNoRows {
		return stats, nil
	}
	return stats, err
}
