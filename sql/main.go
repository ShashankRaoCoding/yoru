package sqlpkg

import (
	"database/sql"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	csvfmt "yoru/sql/csv"
	"yoru/sql/sqldb"
	"yoru/sql/tsv"
	"yoru/utils"
)

type ReaderFunc func(input string) (*sql.DB, error)
type WriterFunc func(db *sql.DB, query string) (string, error)

var readers = map[string]ReaderFunc{
	"csv":   csvfmt.Read,
	"tsv":   tsv.Read,
	"sqldb": sqldb.Read,
}

var writers = map[string]WriterFunc{
	"csv": csvfmt.Write,
	"tsv": tsv.Write,
}

var defaultTableRegex = regexp.MustCompile(`(?i)\b(from|join|update|into|delete from)\s+table\b`)

func Main(args []string) {
	fs := flag.NewFlagSet("sql", flag.ExitOnError)

	inputFormat := fs.String("i", "csv", "Input format: csv, tsv, or sqldb")
	outputFormat := fs.String("o", "csv", "Output format: csv or tsv")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: yoru sql [options] <query>\n")
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		fs.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  cat data.csv | yoru sql -i csv -o csv 'SELECT * FROM table'\n")
		fmt.Fprintf(os.Stderr, "  cat data.tsv | yoru sql -i tsv -o tsv 'SELECT * FROM table'\n")
		fmt.Fprintf(os.Stderr, "  printf 'input.db' | yoru sql -i sqldb -o csv 'SELECT * FROM table'\n")
	}

	err := fs.Parse(args)
	utils.Error(err)

	remainingArgs := fs.Args()
	if len(remainingArgs) < 1 {
		utils.Error(fmt.Errorf("query is required"))
		fs.Usage()
		return
	}

	if len(remainingArgs) > 1 {
		utils.Error(fmt.Errorf("unexpected extra arguments: sqldb path must be provided via stdin"))
		return
	}

	query := remainingArgs[0]

	reader, ok := readers[*inputFormat]
	if !ok {
		utils.Error(fmt.Errorf("unknown input format: %s", *inputFormat))
		return
	}

	writer, ok := writers[*outputFormat]
	if !ok {
		utils.Error(fmt.Errorf("unknown output format: %s", *outputFormat))
		return
	}

	stdinBytes, err := io.ReadAll(os.Stdin)
	utils.Error(err)

	db, err := reader(string(stdinBytes))
	utils.Error(err)
	defer db.Close()

	if *inputFormat == "csv" || *inputFormat == "tsv" {
		query = defaultTableRegex.ReplaceAllString(query, `$1 "table"`)
	}

	output, err := writer(db, query)
	utils.Error(err)

	_, err = fmt.Fprint(os.Stdout, output)
	utils.Error(err)
}
