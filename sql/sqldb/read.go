package sqldb

import (
	"database/sql"
	"fmt"
	"strings"
	"yoru/sql/shared"
)

func Read(input string) (*sql.DB, error) {
	dbPath := strings.TrimSpace(input)
	if dbPath == "" {
		return nil, fmt.Errorf("database path required in stdin when using -i sqldb")
	}

	return shared.OpenDatabaseFile(dbPath)
}
