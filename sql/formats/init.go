package formats

import (
	_ "github.com/mattn/go-sqlite3"
)

// init registers all available format handlers
func init() {
	// Register CSV handlers
	RegisterReader(CSV, &CSVReader{})
	RegisterWriter(CSV, &CSVWriter{})

	// Register TSV handlers
	RegisterReader(TSV, &TSVReader{})
	RegisterWriter(TSV, &TSVWriter{})

	// Register SQLDB handlers
	RegisterReader(SQLDB, &SQLDBReader{})
	RegisterWriter(SQLDB, &SQLDBWriter{})
}
