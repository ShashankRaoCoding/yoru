package tsv

import (
	"database/sql"
	"fmt"
	"io"
	"yoru/sql/shared"
)

func Read(stdin io.Reader) (*sql.DB, error) {
	inputBytes, err := io.ReadAll(stdin)
	if err != nil {
		return nil, fmt.Errorf("reading stdin: %w", err)
	}

	db, err := shared.OpenInMemoryDatabase()
	if err != nil {
		return nil, err
	}

	if err := shared.ReadDelimitedToDB(string(inputBytes), db, "table", '\t'); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
