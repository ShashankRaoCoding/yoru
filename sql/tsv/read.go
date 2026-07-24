package tsv

import (
	"database/sql"
	"yoru/sql/shared"
)

func Read(input string) (*sql.DB, error) {
	db, err := shared.OpenInMemoryDatabase()
	if err != nil {
		return nil, err
	}

	if err := shared.ReadDelimitedToDB(input, db, "table", '\t'); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
