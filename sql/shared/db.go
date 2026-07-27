package shared

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func resolveDBPath(dbPath string) (string, error) {
	if strings.TrimSpace(dbPath) == "" {
		return "", fmt.Errorf("database path is empty")
	}

	absPath, err := filepath.Abs(dbPath)
	if err != nil {
		return "", fmt.Errorf("resolving absolute database path: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(absPath), 0o755); err != nil {
		return "", fmt.Errorf("creating database directory: %w", err)
	}

	return absPath, nil
}

func OpenDatabaseFile(dbPath string) (*sql.DB, error) {
	resolvedPath, err := resolveDBPath(dbPath)
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(resolvedPath); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("checking database file: %w", err)
	}

	db, err := sql.Open("sqlite3", "file:"+resolvedPath+"?cache=shared&mode=rwc")
	if err != nil {
		return nil, fmt.Errorf("opening db: %w", err)
	}

	if _, err := db.Exec("PRAGMA journal_mode=WAL;"); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("validating sqlite database file %q: %w", resolvedPath, err)
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
