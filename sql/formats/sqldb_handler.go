package formats

import (
	"database/sql"
	"fmt"
	"os"
)

// SQLDBReader implements the Reader interface for SQLite database files
type SQLDBReader struct{}

// Read implements Reader.Read for SQLDB
// For SQLDB, the source should be a string path to the database file
func (s *SQLDBReader) Read(source interface{}, db *sql.DB, tableName string) error {
	// For SQLDB input, we don't actually read anything - the "source" IS the database
	// This is a no-op because the database connection is already established
	// from the file path in Main()
	return nil
}

// SQLDBWriter implements the Writer interface for SQLite database files
type SQLDBWriter struct{}

// Write implements Writer.Write for SQLDB
// For SQLDB, the destination should be a string path to the database file
// Since we're already working with an open database, this is handled at a higher level
func (s *SQLDBWriter) Write(db *sql.DB, query string, destination interface{}) error {
	// For SQLDB output, the data stays in the database at the path specified
	// This is a no-op because the database file path is already known from Main()
	return nil
}

// OpenDatabaseFile opens or creates a SQLite database at the given path
func OpenDatabaseFile(dbPath string) (*sql.DB, error) {
	// Verify the file exists or can be created
	if _, err := os.Stat(dbPath); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("checking database file: %w", err)
	}

	// Open the database connection
	db, err := sql.Open("sqlite3", "file:"+dbPath+"?cache=shared&mode=rwc")
	if err != nil {
		return nil, fmt.Errorf("opening db: %w", err)
	}

	return db, nil
}
