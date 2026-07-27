package sqldb

import (
	"fmt"
	"strings"
	"yoru/sql/shared"
)

func Write(sdb shared.Sqldb, query string) (string, error) {
	rows, err := sdb.Query(query)
	if err != nil {
		return "", fmt.Errorf("executing query: %w", err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return "", fmt.Errorf("getting columns: %w", err)
	}

	quotedCols := make([]string, len(cols))
	for i, c := range cols {
		quotedCols[i] = fmt.Sprintf(`"%s" TEXT`, strings.ReplaceAll(c, `"`, `""`))
	}

	var builder strings.Builder
	fmt.Fprintf(&builder, "CREATE TABLE IF NOT EXISTS \"result\" (%s);\n", strings.Join(quotedCols, ", "))

	values := make([]interface{}, len(cols))
	ptrs := make([]interface{}, len(cols))
	for i := range values {
		ptrs[i] = &values[i]
	}

	placeholders := make([]string, len(cols))
	for rows.Next() {
		if err := rows.Scan(ptrs...); err != nil {
			return "", fmt.Errorf("scanning row: %w", err)
		}

		for i, v := range values {
			placeholders[i] = formatSQLValue(v)
		}
		fmt.Fprintf(&builder, "INSERT INTO \"result\" VALUES (%s);\n", strings.Join(placeholders, ", "))
	}

	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("iterating rows: %w", err)
	}

	return builder.String(), nil
}

func formatSQLValue(v interface{}) string {
	if v == nil {
		return "NULL"
	}
	var s string
	switch val := v.(type) {
	case []byte:
		s = string(val)
	default:
		s = fmt.Sprintf("%v", val)
	}
	// Escape single quotes and backslashes for safe SQL string literals.
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, "'", "''")
	return fmt.Sprintf("'%s'", s)
}
