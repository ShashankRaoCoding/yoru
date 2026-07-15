package sqlpkg

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// ImportCSVToSQL reads CSV data from r and stores it in a SQLite table
// named tableName inside dbPath. The first CSV row is treated as the header
// and used for column names (all columns are stored as TEXT).
func ImportCSVToSQL(r io.Reader, dbPath, tableName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "file:"+dbPath+"?cache=shared&mode=rwc")
	if err != nil {
		return nil, fmt.Errorf("opening db: %w", err)
	}

	reader := csv.NewReader(r)

	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("reading header: %w", err)
	}

	// Build CREATE TABLE statement (quote column names to avoid keyword clashes)
	cols := make([]string, len(header))
	for i, h := range header {
		cols[i] = fmt.Sprintf(`"%s" TEXT`, strings.TrimSpace(h))
	}
	createStmt := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%s" (%s)`, tableName, strings.Join(cols, ", "))
	if _, err := db.Exec(createStmt); err != nil {
		return nil, fmt.Errorf("creating table: %w", err)
	}

	// Prepare insert statement
	placeholders := strings.Repeat("?,", len(header))
	placeholders = strings.TrimSuffix(placeholders, ",")
	insertStmt := fmt.Sprintf(`INSERT INTO "%s" VALUES (%s)`, tableName, placeholders)

	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("beginning transaction: %w", err)
	}
	stmt, err := tx.Prepare(insertStmt)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("preparing insert: %w", err)
	}
	defer stmt.Close()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("reading row: %w", err)
		}

		args := make([]interface{}, len(record))
		for i, v := range record {
			args[i] = v
		}
		if _, err := stmt.Exec(args...); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("inserting row: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("committing transaction: %w", err)
	}

	return db, nil
}
