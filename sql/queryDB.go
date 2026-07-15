package sqlpkg

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
)

// QueryToCSV runs the given SQL query against db and writes the result set
// as CSV to w (header row + data rows). Only use with trusted/validated
// query strings — see note below about injection risk.
func QueryToCSV(db *sql.DB, query string, w io.Writer) error {
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

// formatValue converts a scanned column value into its CSV string form,
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
