package formats

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
)

// CSVWriter implements the Writer interface for CSV format
type CSVWriter struct{}

// Write implements Writer.Write for CSV
func (c *CSVWriter) Write(db *sql.DB, query string, destination interface{}) error {
	writer, ok := destination.(io.Writer)
	if !ok {
		return fmt.Errorf("CSV writer expects io.Writer destination")
	}
	return writeDataToDelimited(db, query, writer, ',')
}

// TSVWriter implements the Writer interface for TSV format
type TSVWriter struct{}

// Write implements Writer.Write for TSV
func (t *TSVWriter) Write(db *sql.DB, query string, destination interface{}) error {
	writer, ok := destination.(io.Writer)
	if !ok {
		return fmt.Errorf("TSV writer expects io.Writer destination")
	}
	return writeDataToDelimited(db, query, writer, '\t')
}

// writeDataToDelimited runs a query and writes results as delimited data
func writeDataToDelimited(db *sql.DB, query string, w io.Writer, delimiter rune) error {
	rows, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("executing query: %w", err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("getting columns: %w", err)
	}

	cw := csv.NewWriter(w)
	cw.Comma = delimiter
	defer cw.Flush()

	// Write header
	if err := cw.Write(cols); err != nil {
		return fmt.Errorf("writing header: %w", err)
	}

	// Prepare scan destinations — sql.RawBytes/interface{} handles any column type
	values := make([]interface{}, len(cols))
	ptrs := make([]interface{}, len(cols))
	for i := range values {
		ptrs[i] = &values[i]
	}

	record := make([]string, len(cols))
	for rows.Next() {
		if err := rows.Scan(ptrs...); err != nil {
			return fmt.Errorf("scanning row: %w", err)
		}

		for i, v := range values {
			record[i] = formatValue(v)
		}

		if err := cw.Write(record); err != nil {
			return fmt.Errorf("writing row: %w", err)
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterating rows: %w", err)
	}

	cw.Flush()
	return cw.Error()
}

// formatValue converts a scanned column value into its string form,
// handling NULLs and byte slices (common for TEXT/BLOB columns in SQLite).
func formatValue(v interface{}) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case []byte:
		return string(val)
	default:
		return fmt.Sprintf("%v", val)
	}
}
