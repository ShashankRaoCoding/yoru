package shared

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func OpenDatabaseFile(dbPath string) (*sql.DB, error) {
	if _, err := os.Stat(dbPath); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("checking database file: %w", err)
	}

	db, err := sql.Open("sqlite3", "file:"+dbPath+"?cache=shared&mode=rwc")
	if err != nil {
		return nil, fmt.Errorf("opening db: %w", err)
	}

	return db, nil
}

func OpenInMemoryDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, fmt.Errorf("opening in-memory db: %w", err)
	}
	return db, nil
}
