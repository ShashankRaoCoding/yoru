package sqldb

import (
	"fmt"
	"io"
	"strings"
	"yoru/sql/shared"
)

func Read(stdin io.Reader) (shared.Sqldb, error) {
	pathBytes, err := io.ReadAll(stdin)
	if err != nil {
		return shared.Sqldb{}, fmt.Errorf("reading stdin: %w", err)
	}

	dbPath := strings.TrimSpace(string(pathBytes))
	if dbPath == "" {
		return shared.Sqldb{}, fmt.Errorf("database path required in stdin when using -i sqldb")
	}

	return shared.OpenDatabaseFile(dbPath)
}
