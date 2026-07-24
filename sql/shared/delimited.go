package shared

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

func ReadDelimitedToDB(input string, db *sql.DB, tableName string, delimiter rune) error {
	csvReader := csv.NewReader(strings.NewReader(input))
	csvReader.Comma = delimiter

	header, err := csvReader.Read()
	if err != nil {
		return fmt.Errorf("reading header: %w", err)
	}

	cols := make([]string, len(header))
	for i, h := range header {
		cols[i] = fmt.Sprintf(`"%s" TEXT`, strings.TrimSpace(h))
	}

	createStmt := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%s" (%s)`, tableName, strings.Join(cols, ", "))
	if _, err := db.Exec(createStmt); err != nil {
		return fmt.Errorf("creating table: %w", err)
	}

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
		if err != nil {
			if err == io.EOF {
				break
			}
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

func QueryToDelimitedString(db *sql.DB, query string, delimiter rune) (string, error) {
	rows, err := db.Query(query)
	if err != nil {
		return "", fmt.Errorf("executing query: %w", err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return "", fmt.Errorf("getting columns: %w", err)
	}

	var builder strings.Builder
	writer := csv.NewWriter(&builder)
	writer.Comma = delimiter

	if err := writer.Write(cols); err != nil {
		return "", fmt.Errorf("writing header: %w", err)
	}

	values := make([]interface{}, len(cols))
	ptrs := make([]interface{}, len(cols))
	for i := range values {
		ptrs[i] = &values[i]
	}

	record := make([]string, len(cols))
	for rows.Next() {
		if err := rows.Scan(ptrs...); err != nil {
			return "", fmt.Errorf("scanning row: %w", err)
		}

		for i, v := range values {
			record[i] = formatValue(v)
		}

		if err := writer.Write(record); err != nil {
			return "", fmt.Errorf("writing row: %w", err)
		}
	}

	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("iterating rows: %w", err)
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", err
	}

	return builder.String(), nil
}

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
