package shared

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// Sqldb wraps a *sql.DB and exposes a Query method.
type Sqldb struct {
	db *sql.DB
}

func (s Sqldb) Query(query string) (*sql.Rows, error) {
	return s.db.Query(query)
}

func (s Sqldb) Close() error {
	return s.db.Close()
}

func resolveDBPath(dbPath string) (string, error) {
	if strings.TrimSpace(dbPath) == "" {
		return "", fmt.Errorf("database path is empty")
	}

	// filepath.Abs resolves relative paths from the working directory,
	// which is the correct behavior when a user provides a relative path.
	absPath, err := filepath.Abs(dbPath)
	if err != nil {
		return "", fmt.Errorf("resolving absolute database path: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(absPath), 0o755); err != nil {
		return "", fmt.Errorf("creating database directory: %w", err)
	}

	return absPath, nil
}

func OpenDatabaseFile(dbPath string) (Sqldb, error) {
	resolvedPath, err := resolveDBPath(dbPath)
	if err != nil {
		return Sqldb{}, err
	}

	if _, err := os.Stat(resolvedPath); err != nil && !os.IsNotExist(err) {
		return Sqldb{}, fmt.Errorf("checking database file: %w", err)
	}

	db, err := sql.Open("sqlite3", "file:"+resolvedPath+"?cache=shared&mode=rwc")
	if err != nil {
		return Sqldb{}, fmt.Errorf("opening db: %w", err)
	}

	if _, err := db.Exec("PRAGMA journal_mode=WAL;"); err != nil {
		_ = db.Close()
		return Sqldb{}, fmt.Errorf("validating sqlite database file %q: %w", resolvedPath, err)
	}

	return Sqldb{db: db}, nil
}

func OpenInMemoryDatabase() (Sqldb, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return Sqldb{}, fmt.Errorf("opening in-memory db: %w", err)
	}
	return Sqldb{db: db}, nil
}
