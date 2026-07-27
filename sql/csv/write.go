package csv

import (
	"yoru/sql/shared"
)

func Write(sdb shared.Sqldb, query string) (string, error) {
	return shared.QueryToDelimitedString(sdb, query, ',')
}
