package csv

import (
	"database/sql"
	"fmt"
	"io"
	"yoru/sql/shared"
)

func Read(stdin io.Reader) (*sql.DB, error) {
	input, err := io.ReadAll(stdin)
	if err != nil {
		return nil, fmt.Errorf("reading stdin: %w", err)
	}

	db, err := shared.OpenInMemoryDatabase()
	if err != nil {
		return nil, err
	}

	if err := shared.ReadDelimitedToDB(string(input), db, "table", ','); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
