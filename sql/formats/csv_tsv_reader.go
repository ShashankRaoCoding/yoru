package formats

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

// CSVReader implements the Reader interface for CSV format
type CSVReader struct{}

// Read implements Reader.Read for CSV
func (c *CSVReader) Read(source interface{}, db *sql.DB, tableName string) error {
	reader, ok := source.(io.Reader)
	if !ok {
		return fmt.Errorf("CSV reader expects io.Reader source")
	}
	return readDataToSQL(reader, db, tableName, ',')
}

// TSVReader implements the Reader interface for TSV format
type TSVReader struct{}

// Read implements Reader.Read for TSV
func (t *TSVReader) Read(source interface{}, db *sql.DB, tableName string) error {
	reader, ok := source.(io.Reader)
	if !ok {
		return fmt.Errorf("TSV reader expects io.Reader source")
	}
	return readDataToSQL(reader, db, tableName, '\t')
}

// readDataToSQL is a shared helper for CSV/TSV reading
// It reads delimited data from r and stores it in a SQLite table
func readDataToSQL(r io.Reader, db *sql.DB, tableName string, delimiter rune) error {
	csvReader := csv.NewReader(r)
	csvReader.Comma = delimiter

	header, err := csvReader.Read()
	if err != nil {
		return fmt.Errorf("reading header: %w", err)
	}

	// Build CREATE TABLE statement (quote column names to avoid keyword clashes)
	cols := make([]string, len(header))
	for i, h := range header {
		cols[i] = fmt.Sprintf(`"%s" TEXT`, strings.TrimSpace(h))
	}
	createStmt := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%s" (%s)`, tableName, strings.Join(cols, ", "))
	if _, err := db.Exec(createStmt); err != nil {
		return fmt.Errorf("creating table: %w", err)
	}

	// Prepare insert statement
	placeholders := strings.Repeat("?,", len(header))
	placeholders = strings.TrimSuffix(placeholders, ",")
	insertStmt := fmt.Sprintf(`INSERT INTO "%s" VALUES (%s)`, tableName, placeholders)

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	stmt, err := tx.Prepare(insertStmt)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("preparing insert: %w", err)
	}
	defer stmt.Close()

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("reading row: %w", err)
		}

		args := make([]interface{}, len(record))
		for i, v := range record {
			args[i] = v
		}
		if _, err := stmt.Exec(args...); err != nil {
			tx.Rollback()
			return fmt.Errorf("inserting row: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}
