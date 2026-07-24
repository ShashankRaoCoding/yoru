package csv

import (
	"database/sql"
	"yoru/sql/shared"
)

func Write(db *sql.DB, query string) (string, error) {
	return shared.QueryToDelimitedString(db, query, ',')
}
