package sqlpkg

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strings"
	"yoru/globals"
	"yoru/sql/formats"
	"yoru/utils"
)

func Main(args []string) {
	fs := flag.NewFlagSet("sql", flag.ExitOnError)
	
	inputFormat := fs.String("i", "csv", "Input format: csv, tsv, or sqldb")
	outputFormat := fs.String("o", "csv", "Output format: csv, tsv, or sqldb")
	
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: yoru sql [options] <query> [dbpath for sqldb]\n")
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		fs.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  cat data.csv | yoru sql -i csv -o csv 'SELECT * FROM table'\n")
		fmt.Fprintf(os.Stderr, "  cat data.tsv | yoru sql -i tsv -o tsv 'SELECT * FROM table'\n")
		fmt.Fprintf(os.Stderr, "  yoru sql -i sqldb -o csv 'SELECT * FROM table' input.db\n")
		fmt.Fprintf(os.Stderr, "  cat data.csv | yoru sql -i csv -o sqldb 'SELECT * FROM table' output.db\n")
	}
	
	err := fs.Parse(args)
	utils.Error(err)
	
	remainingArgs := fs.Args()
	if len(remainingArgs) < 1 {
		utils.Error(fmt.Errorf("query is required"))
		fs.Usage()
		return
	}
	
	query := remainingArgs[0]
	
	// Parse input and output formats
	inFmt, err := formats.ParseFormat(*inputFormat)
	utils.Error(err)
	
	outFmt, err := formats.ParseFormat(*outputFormat)
	utils.Error(err)
	
	// Get the reader and writer
	reader, err := formats.GetReader(inFmt)
	utils.Error(err)
	
	writer, err := formats.GetWriter(outFmt)
	utils.Error(err)
	
	var db *sql.DB
	var dbPath string
	
	// Handle input based on format
	switch inFmt {
	case formats.CSV, formats.TSV:
		// For CSV/TSV, create temp database and read from stdin
		tempDBPath := globals.TEMP + "/query.db"
		db, err = formats.OpenDatabaseFile(tempDBPath)
		utils.Error(err)
		defer db.Close()
		defer os.Remove(tempDBPath)
		
		err = reader.Read(os.Stdin, db, "table")
		utils.Error(err)
		
	case formats.SQLDB:
		// For SQLDB input, expect database path as next argument
		if len(remainingArgs) < 2 {
			utils.Error(fmt.Errorf("database path required when using -i sqldb"))
			return
		}
		dbPath = remainingArgs[1]
		
		db, err = formats.OpenDatabaseFile(dbPath)
		utils.Error(err)
		defer db.Close()
		
		err = reader.Read(dbPath, db, "table")
		utils.Error(err)
	}
	
	// Handle output based on format
	switch outFmt {
	case formats.CSV, formats.TSV:
		// For CSV/TSV, write to stdout
		err = writer.Write(db, query, os.Stdout)
		utils.Error(err)
		
	case formats.SQLDB:
		// For SQLDB output, expect database path as additional argument
		outDBPath := ""
		
		// Determine which argument is the output database path
		if inFmt == formats.SQLDB {
			// If input is sqldb, output db path should be at remainingArgs[2]
			if len(remainingArgs) < 3 {
				utils.Error(fmt.Errorf("output database path required when using -o sqldb with -i sqldb"))
				return
			}
			outDBPath = remainingArgs[2]
		} else {
			// If input is csv/tsv, output db path should be at remainingArgs[1]
			if len(remainingArgs) < 2 {
				utils.Error(fmt.Errorf("output database path required when using -o sqldb"))
				return
			}
			outDBPath = remainingArgs[1]
		}
		
		// If output path is different from input, or input is not sqldb, write results
		if inFmt != formats.SQLDB || outDBPath != dbPath {
			outDB, err := formats.OpenDatabaseFile(outDBPath)
			utils.Error(err)
			defer outDB.Close()
			
			// Copy query results to output database
			err = copyQueryResults(db, outDB, query)
			utils.Error(err)
		}
		
		// Write to the output database
		err = writer.Write(db, query, outDBPath)
		utils.Error(err)
	}
}

// copyQueryResults executes a query on sourceDB and inserts results into destDB
func copyQueryResults(sourceDB *sql.DB, destDB *sql.DB, query string) error {
	rows, err := sourceDB.Query(query)
	if err != nil {
		return fmt.Errorf("executing query: %w", err)
	}
	defer rows.Close()
	
	cols, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("getting columns: %w", err)
	}
	
	// Create table in destination
	var colDefs []string
	for _, col := range cols {
		colDefs = append(colDefs, fmt.Sprintf(`"%s" TEXT`, col))
	}
	createStmt := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "results" (%s)`, strings.Join(colDefs, ", "))
	_, err = destDB.Exec(createStmt)
	if err != nil {
		return fmt.Errorf("creating destination table: %w", err)
	}
	
	// Insert results
	placeholders := strings.Repeat("?,", len(cols))
	placeholders = strings.TrimSuffix(placeholders, ",")
	insertStmt := fmt.Sprintf(`INSERT INTO "results" VALUES (%s)`, placeholders)
	
	for rows.Next() {
		values := make([]interface{}, len(cols))
		ptrs := make([]interface{}, len(cols))
		for i := range values {
			ptrs[i] = &values[i]
		}
		
		err := rows.Scan(ptrs...)
		if err != nil {
			return fmt.Errorf("scanning row: %w", err)
		}
		
		_, err = destDB.Exec(insertStmt, values...)
		if err != nil {
			return fmt.Errorf("inserting row: %w", err)
		}
	}
	
	return rows.Err()
}
