package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func OpenSQLite(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func InitSchema(db *sql.DB) error {
	schema := `
		CREATE TABLE IF NOT EXISTS processing_results (
			id TEXT PRIMARY KEY,
			document_id TEXT NOT NULL UNIQUE,
			extracted_text TEXT,
			page_count INTEGER,
			file_hash TEXT UNIQUE,
			is_duplicate BOOLEAN,
			thumbnail_info TEXT,
			processed_at TIMESTAMP,
			error_message TEXT
		);

		CREATE TABLE IF NOT EXISTS document_chunks (
			id TEXT PRIMARY KEY,
			document_id TEXT NOT NULL,
			chunk_index INTEGER NOT NULL,
			content TEXT NOT NULL,
			start_page INTEGER,
			end_page INTEGER,
			char_count INTEGER,
			FOREIGN KEY (document_id) REFERENCES processing_results(document_id)
		);
		
		CREATE INDEX IF NOT EXISTS idx_chunks_doc_id ON document_chunks(document_id);
		CREATE INDEX IF NOT EXISTS idx_chunks_content ON document_chunks(content);
	`
	if _, err := db.Exec(schema); err != nil {
		return err
	}

	return nil
}
