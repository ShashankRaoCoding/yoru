package csv

import (
	"fmt"
	"io"
	"yoru/sql/shared"
)

func Read(stdin io.Reader) (shared.Sqldb, error) {
	inputBytes, err := io.ReadAll(stdin)
	if err != nil {
		return shared.Sqldb{}, fmt.Errorf("reading stdin: %w", err)
	}

	sdb, err := shared.OpenInMemoryDatabase()
	if err != nil {
		return shared.Sqldb{}, err
	}

	if err := shared.ReadDelimitedToDB(string(inputBytes), sdb, "table", ','); err != nil {
		sdb.Close()
		return shared.Sqldb{}, err
	}

	return sdb, nil
}
