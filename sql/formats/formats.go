package formats

import (
	"database/sql"
	"fmt"
	"io"
)

// Format represents a data format that can be read from or written to
type Format string

const (
	CSV   Format = "csv"
	TSV   Format = "tsv"
	SQLDB Format = "sqldb"
)

// Reader defines the interface for reading data into a SQL database
type Reader interface {
	Read(source interface{}, db *sql.DB, tableName string) error
}

// Writer defines the interface for writing data from a SQL database
type Writer interface {
	Write(db *sql.DB, query string, destination interface{}) error
}

// Registry maps format strings to their reader/writer implementations
var (
	readers = make(map[Format]Reader)
	writers = make(map[Format]Writer)
)

// RegisterReader registers a reader for a format
func RegisterReader(f Format, r Reader) {
	readers[f] = r
}

// RegisterWriter registers a writer for a format
func RegisterWriter(f Format, w Writer) {
	writers[f] = w
}

// GetReader returns the reader for a given format, or an error if not found
func GetReader(f Format) (Reader, error) {
	r, ok := readers[f]
	if !ok {
		return nil, fmt.Errorf("no reader registered for format: %s", f)
	}
	return r, nil
}

// GetWriter returns the writer for a given format, or an error if not found
func GetWriter(f Format) (Writer, error) {
	w, ok := writers[f]
	if !ok {
		return nil, fmt.Errorf("no writer registered for format: %s", w)
	}
	return w, nil
}

// ParseFormat converts a string to a Format, returning an error if invalid
func ParseFormat(s string) (Format, error) {
	f := Format(s)
	switch f {
	case CSV, TSV, SQLDB:
		return f, nil
	default:
		return "", fmt.Errorf("unknown format: %s", s)
	}
}
