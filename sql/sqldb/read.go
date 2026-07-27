package sqldb

import (
	"fmt"
	"io"
	"strings"
	"database/sql"
	"yoru/sql/shared"
)

func Read(stdin io.Reader) (*sql.DB, error) {
	pathBytes, err := io.ReadAll(stdin)
	if err != nil {
		return nil, fmt.Errorf("reading stdin: %w", err)
	}

	dbPath := strings.TrimSpace(string(pathBytes))
	if dbPath == "" {
		return nil, fmt.Errorf("database path required in stdin when using -i sqldb")
	}

	return shared.OpenDatabaseFile(dbPath)
}
