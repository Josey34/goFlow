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
	return nil
}
